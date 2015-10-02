// Package implementation controller and URL routing
package caffmux

import (
	"fmt"
	"github.com/tjz101/caffmux/logs"
	"net/http"
)

// Initialize the related, such as log initialize...
func init() {
	CaffLogger = logs.NewLogger()
	err := CaffLogger.SetLogger("console", "")
	if err != nil {
		fmt.Println("init console log error:", err)
	}
	CaffLogger.SetEnableFuncCallDepth(true)
}

type Application struct {
	Handlers *ControllerRegistor
}

func NewApplication() *Application {
	cr := NewControllerRegistor()
	return &Application{cr}
}

// Run the HTTP service, addr for listening on port, such as :8080
func (app *Application) Run(addr string) error {
	CaffLogger.Debug(addr)
	return http.ListenAndServe(addr, app.Handlers)
}

// Configuration controller URL routing rules, support for regular expressions
func (app *Application) Router(path string, c ControllerInterface) *Application {
	CaffLogger.Debug(path)
	app.Handlers.Add(path, c)
	return app
}
