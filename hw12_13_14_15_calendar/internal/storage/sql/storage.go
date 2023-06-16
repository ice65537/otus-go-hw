package dbstore

import (
	"context"
	"fmt"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

/*type Storage interface {
	Init(context.Context, logger.Logger) error
	Upsert(context.Context, logger.Logger, storage.Event) error
	Drop(context.Context, logger.Logger, string) error
	Get(context.Context, logger.Logger, time.Time, time.Time) ([]storage.Event, error)
	Close(context.Context, logger.Logger) error
}*/

type Storage struct {
	db  *sqlx.DB
	dsn string
}

func New(host string, port int, dbname, username, password string) *Storage {
	return &Storage{
		dsn: fmt.Sprintf("host=%s port=%d dbname=%s username=%s password=%s",
			host, port, dbname, username, password,
		),
	}
}

func (s *Storage) Init(ctx context.Context, log logger.Logger) error {
	/*sqlx.DB - обертка над *sql.DB
	sqlx.Tx - обертка над *sql.Tx
	sqlx.Stmt - обертка над *sql.Stmt
	sqlx.NamedStmt - PreparedStatement с поддержкой именованых параметров
	Подключение jmoiron/sqlx :*/
	if db, err := sqlx.Open("pgx", s.dsn); err != nil {
		return err
	} else {
		s.db = db
	}
	return nil
}

func (s *Storage) Close(ctx context.Context, log logger.Logger) error {
	// TODO
	return nil
}

func (s *Storage) Upsert(ctx context.Context, log logger.Logger, evt storage.Event) error {
	/*s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.events[evt.ID]
	header := "Updated event"
	if !ok {
		evt.ID = uuid.New().String()
		header = "Created event"
	}
	s.events[evt.ID] = evt
	log.Info(ctx, "Memstorage.Upsert", fmt.Sprintf("%s %s", header, evt))*/
	return nil
}

func (s *Storage) Drop(ctx context.Context, log logger.Logger, id string) error {
	/*s.mu.Lock()
	defer s.mu.Unlock()
	evt, ok := s.events[id]
	if !ok {
		return log.Error(ctx, "Memstorage.Drop", fmt.Sprintf("unknown event [%s]", id))
	}
	delete(s.events, id)
	log.Info(ctx, "Memstorage.Drop", "Dropped event "+evt.String())*/
	return nil
}

func (s *Storage) Get(ctx context.Context, log logger.Logger, dt1 time.Time, dt2 time.Time,
) ([]storage.Event, error) {
	/*log.Debug(ctx, "Memstorage.Get", fmt.Sprintf("select events from [%s,%s]", dt1, dt2), 1)
	idt1 := dt2int(dt1)
	idt2 := dt2int(dt2)
	if idt2 < idt1 {
		return nil, log.Error(ctx, "Memstorage.Get", fmt.Sprintf("invalid period from %s to %s", dt1, dt2))
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
	log.Debug(ctx, "Memstorage.Get", "%d events selected", 1)*/
	return nil, nil
}
