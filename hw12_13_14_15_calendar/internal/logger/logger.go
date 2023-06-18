package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

type Logger struct {
	AppName string
	Level   string
	Depth   int
	cancel  context.CancelFunc
}

const (
	Fatal   = "FATAL"
	Error   = "ERROR"
	Warning = "WARNING"
	Info    = "INFO"
	Debug   = "DEBUG"
)

func New(appname string, level string, depth int, cf context.CancelFunc) *Logger {
	if appname == "" {
		appname = "Unknown"
	}
	dict := map[string]struct{}{
		Fatal:   {},
		Error:   {},
		Warning: {},
		Info:    {},
		Debug:   {},
	}
	if _, ok := dict[level]; !ok {
		level = Error
	}
	switch {
	case level != Debug:
		depth = 0
	case depth < 1:
		depth = 1
	case depth > 5:
		depth = 5
	}
	return &Logger{AppName: appname, Level: level, Depth: depth, cancel: cf}
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
		return fmt.Sprintf(`{"timestamp"="%s","level"=Error,"depth"="0","oper"="Logger.Encode","text"="%v"}`,
			time.Now(), err)
	}
	return string(rslt)
}

func (l Logger) output(ctx context.Context, oper, txt, level string, depth int) {
	f := os.Stdout
	if level == Fatal || level == Error || level == Warning {
		f = os.Stderr
	}
	fmt.Fprintln(f, l.encode(ctx, oper, txt, level, depth))
}

func (l Logger) Fatal(ctx context.Context, oper string, msg any) error {
	defer l.cancel()
	l.output(ctx, oper, fmt.Sprintf("%v", msg), Fatal, 0)
	return fmt.Errorf(strings.ToLower(oper)+": %v", msg)
}

func (l Logger) Error(ctx context.Context, oper string, msg any) error {
	if l.Level != Fatal {
		l.output(ctx, oper, fmt.Sprintf("%v", msg), Error, 0)
	}
	return fmt.Errorf(strings.ToLower(oper)+": %v", msg)
}

func (l Logger) Warning(ctx context.Context, oper, msg string) {
	if l.Level == Error || l.Level == Fatal {
		return
	}
	l.output(ctx, oper, msg, Warning, 0)
}

func (l Logger) Info(ctx context.Context, oper, msg string) {
	if l.Level == Error || l.Level == Warning || l.Level == Fatal {
		return
	}
	l.output(ctx, oper, msg, Info, 0)
}

func (l Logger) Debug(ctx context.Context, oper, msg string, depth int) {
	if l.Level == Error || l.Level == Warning || l.Level == Info ||
		l.Level == Fatal || depth > l.Depth {
		return
	}
	if depth < 1 {
		depth = 1
	}
	l.output(ctx, oper, msg, Debug, depth)
}
