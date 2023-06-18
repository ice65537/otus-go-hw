package memstore

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
	log    *logger.Logger
}

func New() *Storage {
	return &Storage{events: map[string]storage.Event{}}
}

func (s *Storage) Init(_ context.Context, log *logger.Logger, _ string) error {
	s.log = log
	return nil
}

func (s *Storage) Upsert(ctx context.Context, evt storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	header := "Updated event"
	if _, ok := s.events[evt.ID]; !ok {
		evt.ID = uuid.New().String()
		header = "Created event"
	}
	s.events[evt.ID] = evt
	s.log.Info(ctx, "Memstorage.Upsert", fmt.Sprintf("%s %s", header, evt))
	return nil
}

func (s *Storage) Drop(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	evt, ok := s.events[id]
	if !ok {
		return s.log.Error(ctx, "Memstorage.Drop", fmt.Sprintf("unknown event [%s]", id))
	}
	delete(s.events, id)
	s.log.Info(ctx, "Memstorage.Drop", "Dropped event "+evt.String())
	return nil
}

func (s *Storage) Get(ctx context.Context, dt1 time.Time, dt2 time.Time,
) ([]storage.Event, error) {
	//
	s.mu.RLock()
	defer s.mu.RUnlock()
	//
	result := []storage.Event{}
	for _, v := range s.events {
		if (v.StartDt.After(dt1) || v.StartDt.Equal(dt1)) &&
			(v.StartDt.Before(dt2) || v.StartDt.Equal(dt2)) {
			result = append(result, v)
		}
	}
	s.log.Debug(ctx, "Memstorage.Get", "%d events selected", 1)
	return result, nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}
