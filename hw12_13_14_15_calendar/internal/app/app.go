package app

import (
	"context"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

type App struct {
	log   *logger.Logger
	store Storage
}

type Storage interface { // TODO
}

func New(logLevel string, logDepth int, store Storage) *App {
	app := App{}
	app.log = logger.New(logLevel, logDepth)
	app.store = store
	return &app
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
