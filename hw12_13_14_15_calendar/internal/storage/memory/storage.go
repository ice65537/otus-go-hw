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

func (s *Storage) Init(ctx context.Context, log *logger.Logger) error {
	s.log = log
	return nil
}

func (s *Storage) Upsert(ctx context.Context, evt storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.events[evt.ID]
	header := "Updated event"
	if !ok {
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
	s.log.Debug(ctx, "Memstorage.Get", fmt.Sprintf("select events from [%s,%s]", dt1, dt2), 1)
	idt1 := dt2int(dt1)
	idt2 := dt2int(dt2)
	if idt2 < idt1 {
		return nil, s.log.Error(ctx, "Memstorage.Get", fmt.Sprintf("invalid period from %s to %s", dt1, dt2))
	}
	//
	s.mu.RLock()
	defer s.mu.RUnlock()
	//
	result := []storage.Event{}
	for _, v := range s.events {
		if dt2int(v.StartDT) >= idt1 && dt2int(v.StartDT) <= idt2 {
			result = append(result, v)
		}
	}
	s.log.Debug(ctx, "Memstorage.Get", "%d events selected", 1)
	return result, nil
}

func dt2int(dateTime time.Time) int {
	return dateTime.Year()*1000 + dateTime.YearDay()
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}
