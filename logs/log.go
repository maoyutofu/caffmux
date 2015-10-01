// Package has realized the log and adaptation
package logs

import (
	"fmt"
	"path"
	"runtime"
)

const (
	LevelError = iota
	LevelWarning
	LevelInfo
	LevelDebug
)

type loggerType func() LoggerInterface

type LoggerInterface interface {
	Init(l string) error
	WriteMsg(msg string, level int) error
	Destory()
}

var adapters = make(map[string]loggerType)

// According to the different log output mode of adapter name to register
func Register(name string, log loggerType) {
	if log == nil {
		panic("logs: Register provide is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("logs: Register called twice for provider " + name)
	}
	adapters[name] = log
}

type CaffLogger struct {
	level               int
	enableFuncCallDepth bool
	loggerFuncCallDepth int
	outputs             map[string]LoggerInterface
}

// Instantiation CaffLogger structure, and return
func NewLogger() *CaffLogger {
	cl := new(CaffLogger)
	cl.level = LevelDebug
	cl.loggerFuncCallDepth = 2
	cl.outputs = make(map[string]LoggerInterface)
	return cl
}

// Set the log level
func (cl *CaffLogger) SetLevel(l int) {
	cl.level = l
}

func (cl *CaffLogger) SetLoggerFuncCallDepth(d int) {
	cl.loggerFuncCallDepth = d
}

func (cl *CaffLogger) GetLoggerFuncCallDepth() int {
	return cl.loggerFuncCallDepth
}

func (cl *CaffLogger) SetEnableFuncCallDepth(b bool) {
	cl.enableFuncCallDepth = b
}

// By specifying the name of the adapter to initialize a log object
func (cl *CaffLogger) SetLogger(adapter string, jsonconf string) error {
	if log, ok := adapters[adapter]; ok {
		lg := log()
		err := lg.Init(jsonconf)
		cl.outputs[adapter] = lg
		if err != nil {
			fmt.Println("logs.CaffLogger.SetLogger: " + err.Error())
			return err
		}
	} else {
		return fmt.Errorf("logs: unknown adapter %q (forgotten Register?)", adapter)
	}
	return nil
}

// According to the adapter name removed from the log object
func (cl *CaffLogger) DeleteLogger(adapter string) error {
	if lg, ok := cl.outputs[adapter]; ok {
		lg.Destory()
		delete(cl.outputs, adapter)
		return nil
	} else {
		return fmt.Errorf("logs: unknown adapter %q (forgotten Register?)", adapter)
	}
}

func (cl *CaffLogger) write(level int, msg string) error {
	if cl.enableFuncCallDepth {
		_, file, line, ok := runtime.Caller(cl.loggerFuncCallDepth)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		msg = fmt.Sprintf("[%s:%d] %s", filename, line, msg)
	}
	for name, l := range cl.outputs {
		err := l.WriteMsg(msg, level)
		if err != nil {
			fmt.Println("unable to write to adapter:", name, err)
			return err
		}
	}
	return nil
}

// Specifies the format of the output error log
func (cl *CaffLogger) Error(format string, v ...interface{}) {
	if LevelError > cl.level {
		return
	}
	msg := fmt.Sprintf("[C] "+format, v...)
	cl.write(LevelError, msg)
}

// Specifies the format of the output warning log
func (cl *CaffLogger) Warning(format string, v ...interface{}) {
	if LevelWarning > cl.level {
		return
	}
	msg := fmt.Sprintf("[W] "+format, v...)
	cl.write(LevelWarning, msg)
}

// Specifies the format of the output info log
func (cl *CaffLogger) Info(format string, v ...interface{}) {
	if LevelInfo > cl.level {
		return
	}
	msg := fmt.Sprintf("[I] "+format, v...)
	cl.write(LevelInfo, msg)
}

// Specifies the format of the output debug log
func (cl *CaffLogger) Debug(format string, v ...interface{}) {
	if LevelDebug > cl.level {
		return
	}
	msg := fmt.Sprintf("[D] "+format, v...)
	cl.write(LevelDebug, msg)
}
