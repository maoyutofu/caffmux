package logs

import (
	"encoding/json"
	"log"
	"os"
)

type ConsoleWrite struct {
	lg    *log.Logger
	Level int `json:"level"`
}

func NewConsole() LoggerInterface {
	cw := &ConsoleWrite{
		lg:    log.New(os.Stdout, "", log.Ldate|log.Ltime),
		Level: LevelDebug,
	}
	return cw
}

func (c *ConsoleWrite) Init(jsonconf string) error {
	if len(jsonconf) == 0 {
		return nil
	}
	return json.Unmarshal([]byte(jsonconf), c)
}

func (c *ConsoleWrite) Write(msg string, level int) error {
	if level > c.Level {
		return nil
	}
	c.lg.Println(msg)
	return nil
}

func (c *ConsoleWrite) Destory() {

}

func init() {
	Register("console", NewConsole)
}
