package logs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type FileWriter struct {
	lg          *log.Logger
	o           *FileOut
	Level       int    `json:"level"`
	Filename    string `json:"filename"`
	logFilename string
	ts          string
	expression  string
}

type FileOut struct {
	sync.Mutex
	fp *os.File
}

func (o *FileOut) Write(b []byte) (int, error) {
	o.Lock()
	defer o.Unlock()
	return o.fp.Write(b)
}

func (o *FileOut) SetFp(fp *os.File) {
	if o.fp != nil {
		o.fp.Close()
	}
	o.fp = fp
}

func NewFile() LoggerInterface {
	fw := &FileWriter{
		Filename: "",
		Level:    LevelDebug,
	}
	fw.o = new(FileOut)
	fw.lg = log.New(fw.o, "", log.Ldate|log.Ltime)
	return fw
}

func (w *FileWriter) Init(jsonconf string) error {
	err := json.Unmarshal([]byte(jsonconf), w)
	if err != nil {
		return err
	}
	if len(w.Filename) == 0 {
		return errors.New("jsonconfig must have filename")
	}
	err = w.startLogger()
	return err
}

func (w *FileWriter) startLogger() error {
	fp, err := w.createLogFile()
	if err != nil {
		return err
	}
	w.o.SetFp(fp)
	return nil
}

func extract(s string) string {
	var expr string
	r := regexp.MustCompile("\\((\\S+)\\)")
	if r.MatchString(s) {
		expr = r.FindStringSubmatch(s)[1]
	}
	return expr
}

func (w *FileWriter) createLogFile() (*os.File, error) {
	w.expression = extract(w.Filename)
	if w.expression != "" {
		sep := fmt.Sprintf("(%s)", w.expression)
		names := strings.SplitN(w.Filename, sep, 2)
		t := time.Now()
		format := t.Format(w.expression)
		w.ts = format
		w.logFilename = fmt.Sprintf("%s%s%s", names[0], format, names[1])
	}
	fp, err := os.OpenFile(w.logFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	return fp, err
}

func (w *FileWriter) isToday() bool {
	t := time.Now()
	format := t.Format(w.expression)
	return w.ts == format
}

func (w *FileWriter) WriteMsg(msg string, level int) error {
	if level > w.Level {
		return nil
	}
	if !w.isToday() {
		w.Destory()
		w.startLogger()
	}
	w.lg.Println(msg)
	return nil
}

func (w *FileWriter) Destory() {
	w.o.fp.Close()
}

func (w *FileWriter) Flush() {
	w.o.fp.Sync()
}

func init() {
	Register("file", NewFile)
}
