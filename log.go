package caffmux

import (
	"github.com/caffmux/logs"
	"strings"
)

const (
	LevelError = iota
	LevelWarning
	LevelInfo
	LevelDebug
)

var CaffLogger *logs.CaffLogger

func SetLevel(l int) {
	CaffLogger.SetLevel(l)
}

func SetLogger(adapter string, jsonconf string) error {
	err := CaffLogger.SetLogger(adapter, jsonconf)
	if err != nil {
		return err
	}
	return nil
}

func Error(v ...interface{}) {
	CaffLogger.Error(generateFormatString(len(v)), v...)
}

func Warning(v ...interface{}) {
	CaffLogger.Warning(generateFormatString(len(v)), v...)
}

func Info(v ...interface{}) {
	CaffLogger.Info(generateFormatString(len(v)), v...)
}

func Debug(v ...interface{}) {
	CaffLogger.Debug(generateFormatString(len(v)), v...)
}

func generateFormatString(n int) string {
	return strings.Repeat("%v ", n)
}
