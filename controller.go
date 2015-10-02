package caffmux

import (
	"net/http"
)

type Controller struct {
	cname   string
	Content *Context
}

type ControllerInterface interface {
	Init(cname string, content *Context)
	Prepare()
	Get()
	Post()
	Delete()
	Put()
	Head()
	Patch()
	Options()
	Finish()
	Render() error
}

func (c *Controller) Init(cname string, content *Context) {
	c.cname = cname
	c.Content = content
}

func (c *Controller) Prepare() {

}
func (c *Controller) Get() {
	http.Error(c.Content.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Post() {
	http.Error(c.Content.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Delete() {
	http.Error(c.Content.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Put() {
	http.Error(c.Content.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Head() {
	http.Error(c.Content.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Patch() {
	http.Error(c.Content.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Options() {
	http.Error(c.Content.ResponseWriter, "Method Not Allowed", 405)
}
func (c *Controller) Render() error {
	return nil
}
func (c *Controller) Finish() {

}
