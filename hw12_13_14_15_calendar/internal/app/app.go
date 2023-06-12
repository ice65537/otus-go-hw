package app

import (
	"context"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Log   *logger.Logger
	store Storage
}

type Storage interface {
	Init(context.Context, logger.Logger) error
	Upsert(context.Context, logger.Logger, storage.Event) error
	Drop(context.Context, logger.Logger, string) error
	Get(context.Context, logger.Logger, time.Time, time.Time) ([]storage.Event, error)
	Close(context.Context, logger.Logger) error
}

func New(appName, logLevel string, logDepth int, store Storage) *App {
	app := App{}
	app.Log = logger.New(appName, logLevel, logDepth)
	app.store = store
	return &app
}

/*func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}*/
// TODO
