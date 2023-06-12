package memstore

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{events: map[string]storage.Event{}}
}

func (s *Storage) Init(ctx context.Context) error {
	return nil
}

func (s *Storage) Upsert(evt storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_,ok:=s.events[evt.ID]
	if !ok {
		evt.ID = uuid.New().String()
	}
	s.events[evt.ID]=evt
	return nil
}

func (s *Storage) Drop(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_,ok:=s.events[id]
	if !ok {
		return fmt.Errorf("unknown event [%s]",id)
	}
}

// TODO
/*

Get(time.Time, time.Time) ([]storage.Event, error)
Close(ctx context.Context) error*/
