package caffmux

import (
	"caffmux/logs"
	"fmt"
	"net/http"
)

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

func (app *Application) Run(addr string) error {
	CaffLogger.Debug(addr)
	return http.ListenAndServe(addr, app.Handlers)
}

func (app *Application) Router(path string, c ControllerInterface) *Application {
	CaffLogger.Debug(path)
	app.Handlers.Add(path, c)
	return app
}
