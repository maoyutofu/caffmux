package caffmux

import (
	"fmt"
	"github.com/tjz101/caffmux/logs"
)

var (
	StaticPath map[string]string
)

// Initialize the related, such as log initialize...
func init() {
	CaffLogger = logs.NewLogger()
	//err := CaffLogger.SetLogger("file", `{"filename":"localhost.(2006-01-02).log"}`)
	err := CaffLogger.SetLogger("console", "")
	if err != nil {
		fmt.Println("init console log error:", err)
	}
	SetLogFuncCall(true)
	StaticPath = make(map[string]string)
}
