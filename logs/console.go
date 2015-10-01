package logs

import (
	"encoding/json"
	"log"
	"os"
)

type ConsoleWriter struct {
	lg    *log.Logger
	Level int `json:"level"`
}

func NewConsole() LoggerInterface {
	cw := &ConsoleWriter{
		lg:    log.New(os.Stdout, "", log.Ldate|log.Ltime),
		Level: LevelDebug,
	}
	return cw
}

func (c *ConsoleWriter) Init(jsonconf string) error {
	if len(jsonconf) == 0 {
		return nil
	}
	return json.Unmarshal([]byte(jsonconf), c)
}

func (c *ConsoleWriter) WriteMsg(msg string, level int) error {
	if level > c.Level {
		return nil
	}
	c.lg.Println(msg)
	return nil
}

func (c *ConsoleWriter) Destory() {

}

func init() {
	Register("console", NewConsole)
}
