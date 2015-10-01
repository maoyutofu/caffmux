// Package implementation controller and URL routing
package caffmux

import (
	"net/http"
)

type Application struct {
	Handlers *ControllerRegistor
}

func NewApplication() *Application {
	cr := NewControllerRegistor()
	return &Application{cr}
}

// To add a static route
func (app *Application) AddStaticPath(url string, path string) *Application {
	StaticPath[url] = path
	return app
}

// Run the HTTP service, addr for listening on port, such as :8080
func (app *Application) Run(addr string) error {
	Debug("http server listener " + addr)
	return http.ListenAndServe(addr, app.Handlers)
}

// Configuration controller URL routing rules, support for regular expressions
func (app *Application) Router(path string, c ControllerInterface) *Application {
	Debug(path)
	app.Handlers.Add(path, c)
	return app
}
