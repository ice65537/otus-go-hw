package server

import (
	"context"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Application interface {
	Logger() *logger.Logger
	Upsert(context.Context, storage.Event) error
	Drop(context.Context, string) error
	Get(context.Context, time.Time, time.Time) ([]storage.Event, error)
}
