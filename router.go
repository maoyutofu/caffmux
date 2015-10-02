package caffmux

import (
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

type ControllerInfo struct {
	regex  *regexp.Regexp
	params map[int]string
	cType  reflect.Type
}

type ControllerRegistor struct {
	routers []*ControllerInfo
}

func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{}
}

func (cr *ControllerRegistor) Add(pattern string, c ControllerInterface) {
	parts := strings.Split(pattern, "/")
	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"
			if index := strings.Index(part, "("); index >= 3 {
				expr = part[index:]
				part = part[1:index]
			}
			params[j] = part
			parts[i] = expr
			j++
		}
	}
	pattern = strings.Join(parts, "/")
	CaffLogger.Debug(pattern)
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		panic(regexErr)
		return
	}
	t := reflect.Indirect(reflect.ValueOf(c)).Type()
	route := &ControllerInfo{}
	route.regex = regex
	route.params = params
	route.cType = t
	cr.routers = append(cr.routers, route)
}

func (cr *ControllerRegistor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "HTTP File Storage")
	w.Header().Set("Author", "Tang Jizhong")
	var started bool
	requestPath := r.URL.Path
	CaffLogger.Debug(requestPath)
	for _, route := range cr.routers {
		if !route.regex.MatchString(requestPath) {
			continue
		}
		matches := route.regex.FindStringSubmatch(requestPath)
		if len(matches[0]) != len(requestPath) {
			continue
		}
		params := make(map[string]string)
		if len(route.params) > 0 {
			values := r.URL.Query()
			for i, match := range matches[1:] {
				values.Add(route.params[i], match)
				params[route.params[i]] = match
			}
			r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
		}
		c := reflect.New(route.cType)
		init := c.MethodByName("Init")
		in := make([]reflect.Value, 2)
		in[0] = reflect.ValueOf(route.cType.Name())
		in[1] = reflect.ValueOf(&Context{r, w, params})
		init.Call(in)
		in = make([]reflect.Value, 0)
		actionName := ""
		actionNames := strings.Split(requestPath, "!")
		if len(actionNames) > 1 && actionNames[1] != "" {
			actionName = actionNames[1]
		} else {
			actionNames = strings.Split(requestPath, "/")
			if len(actionNames) > 1 && actionNames[1] != "" {
				actionName = actionNames[1]
			}
		}
		if actionName != "" {
			CaffLogger.Debug(actionName)
			numMethod := c.NumMethod()
			t := c.Type()
			for i := 0; i < numMethod; i++ {
				CaffLogger.Debug(t.Method(i).Name)
				if strings.ToUpper(t.Method(i).Name) == strings.ToUpper(actionName) {
					c.Method(i).Call(in)
					break
				}
			}
		} else {
			method := c.MethodByName("Prepare")
			method.Call(in)
			if r.Method == "GET" {
				method = c.MethodByName("Get")
				method.Call(in)
			} else if r.Method == "POST" {
				method = c.MethodByName("Post")
				method.Call(in)
			} else if r.Method == "HEAD" {
				method = c.MethodByName("Head")
				method.Call(in)
			} else if r.Method == "DELETE" {
				method = c.MethodByName("Delete")
				method.Call(in)
			} else if r.Method == "PUT" {
				method = c.MethodByName("Put")
				method.Call(in)
			} else if r.Method == "PATCH" {
				method = c.MethodByName("Patch")
				method.Call(in)
			} else if r.Method == "OPTIONS" {
				method = c.MethodByName("Options")
				method.Call(in)
			}
			method = c.MethodByName("Finish")
			method.Call(in)
		}
		started = true
	}
	if !started {
		http.NotFound(w, r)
	}
}