package app

import (
	"context"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	log   *logger.Logger
	store Storage
}

type Storage interface {
	Init(context.Context, *logger.Logger) error
	Upsert(context.Context, storage.Event) error
	Drop(context.Context, string) error
	Get(context.Context, time.Time, time.Time) ([]storage.Event, error)
	Close(context.Context) error
}

func New(appName, logLevel string, logDepth int, store Storage) *App {
	app := App{}
	app.log = logger.New(appName, logLevel, logDepth)
	app.store = store
	return &app
}

func (a App) Logger() *logger.Logger {
	return a.log
}
