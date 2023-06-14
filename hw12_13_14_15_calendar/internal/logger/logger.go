package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
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

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
	AppName   string    `json:"appname"`
	Level     string    `json:"level"`
	Depth     int       `json:"depth"`
	Oper      string    `json:"oper"`
	Session   Session   `json:"session,omitempty"`
}
type Session struct {
	UUID  string    `json:"uuid"`
	User  string    `json:"user,omitempty"`
	Start time.Time `json:"start"`
}
type logCtxKey string

const keyCtxSession logCtxKey = "keyCtxSession"

func GetCtxSession(ctx context.Context) Session {
	value, ok := ctx.Value(keyCtxSession).(Session)
	if !ok {
		value = Session{User: "", UUID: uuid.New().String(), Start: time.Now()}
	}
	return value
}

func PushCtxSession(ctx context.Context, obj Session) context.Context {
	return context.WithValue(ctx, keyCtxSession, obj)
}

func (l Logger) encode(ctx context.Context, oper, txt, level string, depth int) string {
	msg := Message{
		Timestamp: time.Now(), Text: txt, AppName: l.AppName,
		Level: level, Depth: depth, Oper: oper, Session: GetCtxSession(ctx),
	}
	rslt, err := json.Marshal(msg)
	if err != nil {
		return fmt.Sprintf(`{"timestamp"="%s","level"="ERROR","depth"="0","oper"="Logger.Encode","text"="%v"}`,
			time.Now(), err)
	}
	return string(rslt)
}

func (l Logger) output(ctx context.Context, oper, txt, level string, depth int) {
	f := os.Stdout
	if level == "ERROR" || level == "WARNING" {
		f = os.Stderr
	}
	fmt.Fprintln(f, l.encode(ctx, oper, txt, level, depth))
}

func (l Logger) ErrorE(ctx context.Context, oper string, err error) error {
	l.output(ctx, oper, fmt.Sprintf("%v", err), "ERROR", 0)
	return err
}

func (l Logger) Error(ctx context.Context, oper, msg string) error {
	l.output(ctx, oper, msg, "ERROR", 0)
	return fmt.Errorf(strings.ToLower(oper) + ": " + msg)
}

func (l Logger) Warning(ctx context.Context, oper, msg string) {
	if l.Level == "ERROR" {
		return
	}
	l.output(ctx, oper, msg, "WARNING", 0)
}

func (l Logger) Info(ctx context.Context, oper, msg string) {
	if l.Level == "ERROR" || l.Level == "WARNING" {
		return
	}
	l.output(ctx, oper, msg, "INFO", 0)
}

func (l Logger) Debug(ctx context.Context, oper, msg string, depth int) {
	if l.Level == "ERROR" || l.Level == "WARNING" || l.Level == "INFO" || depth > l.Depth {
		return
	}
	if depth < 1 {
		depth = 1
	}
	l.output(ctx, oper, msg, "DEBUG", depth)
}
