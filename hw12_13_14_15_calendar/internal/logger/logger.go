package logger

import (
	"fmt"
	"time"
)

type Logger struct {
	AppName string
	Level   string
	Depth   int
}

func New(appname string, level string, depth int) *Logger {
	if appname == "" {
		appname = "Unknown"
	}
	dict := map[string]struct{}{
		"ERROR":   {},
		"WARNING": {},
		"INFO":    {},
		"DEBUG":   {},
	}
	if _, ok := dict[level]; !ok {
		level = "ERROR"
	}
	if level != "DEBUG" {
		depth = 0
	} else if depth < 1 {
		depth = 1
	} else if depth > 5 {
		depth = 5
	}
	return &Logger{AppName: appname, Level: level, Depth: depth}
}

/*
type Message struct {
	Timestamp time.Time `json:timestamp`
	AppName   string    `json:appname`
	Level     string    `json:level`
	Depth     int       `json:depth`
	Oper      string    `json:oper`
	Text      string    `json:text`
}
*/

func (l Logger) encode(oper, txt, level string, depth int) string {
	fmtS := `{"timestamp"="%s";"appname"="%s";"level"="%s";"depth"="%d";"oper"="%s";"text"="%s"}` //nolint:gocritic
	return fmt.Sprintf(fmtS, time.Now().Format("RFC3339"), l.AppName, level, depth, oper, txt)
}

func (l Logger) output(oper, txt, level string, depth int) {
	fmt.Println(l.encode(oper, txt, level, depth))
}

func (l Logger) Error(oper, msg string) {
	l.output(oper, msg, "ERROR", 0)
}

func (l Logger) Warning(oper, msg string) {
	if l.Level == "ERROR" {
		return
	}
	l.output(oper, msg, "WARNING", 0)
}

func (l Logger) Info(oper, msg string) {
	if l.Level == "ERROR" || l.Level == "WARNING" {
		return
	}
	l.output(oper, msg, "INFO", 0)
}

func (l Logger) Debug(oper, msg string, depth int) {
	if l.Level == "ERROR" || l.Level == "WARNING" || l.Level == "INFO" || depth > l.Depth {
		return
	}
	if depth < 1 {
		depth = 1
	}
	l.output(oper, msg, "DEBUG", depth)
}
