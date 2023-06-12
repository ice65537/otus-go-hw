package logger

import "fmt"

type Logger struct {
	Level string
	Depth int
}

func New(level string, depth int) *Logger {
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
	return &Logger{Level: level, Depth: depth}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	// TODO
}

// TODO
