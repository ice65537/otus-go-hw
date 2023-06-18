package app

import (
	"context"
	"fmt"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	log    *logger.Logger
	store  *Storage
	cancel context.CancelFunc
}

type Storage interface {
	Init(context.Context, *logger.Logger, string) error
	Upsert(context.Context, storage.Event) error
	Drop(context.Context, string) error
	Get(context.Context, time.Time, time.Time) ([]storage.Event, error)
	Close(context.Context) error
}

func New(appName,
	logLevel string, logDepth int,
	store Storage,
	cf context.CancelFunc,
) *App {
	app := App{}
	app.log = logger.New(appName, logLevel, logDepth, cf)
	app.store = &store
	app.cancel = cf
	return &app
}

func (a App) Logger() *logger.Logger {
	return a.log
}

func (a App) Init(ctx context.Context, connStr string) error {
	return (*a.store).Init(ctx, a.log, connStr)
}

func (a App) Upsert(ctx context.Context, evt storage.Event) error {
	return (*a.store).Upsert(ctx, evt)
}

func (a App) Drop(ctx context.Context, id string) error {
	return (*a.store).Drop(ctx, id)
}

func (a App) Get(ctx context.Context, t1 time.Time, t2 time.Time) ([]storage.Event, error) {
	a.log.Debug(ctx, "App.Get", fmt.Sprintf("select events from [%s,%s]", t1, t2), 1)
	idt1 := storage.Dt2int(t1)
	idt2 := storage.Dt2int(t2)
	if idt2 < idt1 {
		return nil, a.log.Error(ctx, "App.Get", fmt.Sprintf("invalid period from %s to %s", t1, t2))
	}
	return (*a.store).Get(ctx, t1, t2)
}

func (a App) Close(ctx context.Context) error {
	defer a.cancel()
	return (*a.store).Close(ctx)
}
