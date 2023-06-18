package dbstore

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db     *sqlx.DB
	log    *logger.Logger
	insert *sqlx.NamedStmt
	update *sqlx.NamedStmt
	delete *sqlx.NamedStmt
	get    *sqlx.NamedStmt
	getOne *sqlx.NamedStmt
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Init(ctx context.Context, log *logger.Logger, connStr string) error {
	var err error
	s.log = log
	if s.db, err = sqlx.Open("pgx", connStr); err != nil {
		return s.log.Fatal(ctx, "DB.Connect", err)
	}

	x := []int64{}
	err = s.db.Select(&x, "select 1 as x")
	if err != nil {
		return s.log.Fatal(ctx, "DB.ConnChk", err)
	}

	s.insert, err = s.db.PrepareNamedContext(ctx,
		`insert into t_event (eid,etitle,estartdt,estopdt,edesc,eowner,enotifybefore)
		values(:id,:title,:startdt,:stopdt,:desc,:owner,:notifybefore)`,
	)
	if err != nil {
		s.log.Fatal(ctx, "DB.InsPrepare", err)
		return err
	}
	s.update, err = s.db.PrepareNamedContext(ctx,
		`update t_event set 
		etitle=:title,estartdt=:startdt,estopdt=:stopdt,
		edesc=:desc,eowner=:owner,enotifybefore=:notifybefore		
		where eid=:id`,
	)
	if err != nil {
		s.log.Fatal(ctx, "DB.UpdPrepare", err)
		return err
	}
	s.delete, err = s.db.PrepareNamedContext(ctx, `delete from t_event where eid=:id`)
	if err != nil {
		s.log.Fatal(ctx, "DB.DelPrepare", err)
		return err
	}
	s.get, err = s.db.PrepareNamedContext(ctx,
		`select eid,etitle,estartdt,estopdt,eowner,enotifybefore,coalesce(edesc,'') 
		from t_event 
		where estartdt between :dt1 and :dt2`)
	if err != nil {
		s.log.Fatal(ctx, "DB.GetPrepare", err)
		return err
	}
	s.getOne, err = s.db.PrepareNamedContext(ctx,
		`select eid,etitle,estartdt,estopdt,eowner,enotifybefore,coalesce(edesc,'') 
		from t_event 
		where eid = :id`)
	if err != nil {
		s.log.Fatal(ctx, "DB.GetOnePrepare", err)
		return err
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if err := s.db.Close(); err != nil {
		return s.log.Error(ctx, "DB.Close", err)
	}
	return nil
}

func (s *Storage) Upsert(ctx context.Context, evt storage.Event) error {
	evt.StartDt = evt.StartDt.Truncate(time.Minute)
	evt.StopDt = evt.StopDt.Truncate(time.Minute)
	header := "Update"
	exe := s.update
	if evt.ID == "" {
		evt.ID = uuid.New().String()
		header = "Create"
		exe = s.insert
	}
	sqlr, err := exe.ExecContext(ctx,
		map[string]any{
			"id":           evt.ID,
			"title":        evt.Title,
			"startdt":      evt.StartDt,
			"stopdt":       evt.StopDt,
			"owner":        evt.Owner,
			"notifybefore": evt.NotifyBefore,
			"desc":         evt.Desc,
		})
	if err != nil {
		return s.log.Error(ctx, "DB.Upsert."+header, err)
	}
	if header == "Update" {
		n, err := sqlr.RowsAffected()
		if err == nil && n != 1 {
			return s.log.Error(ctx, "DB.Update",
				fmt.Sprintf("unexpected rows affected count [%d]", n))
		}
	}
	s.log.Info(ctx, "DB.Upsert."+header, evt.String())
	return nil
}

func (s *Storage) Drop(ctx context.Context, id string) error {
	row := s.getOne.QueryRowContext(ctx, map[string]any{"id": id})
	if row.Err() != nil {
		return s.log.Error(ctx, "DB.DropChkGet", row.Err())
	}
	evt := storage.Event{}
	if err := row.Scan(&evt.ID, &evt.Title, &evt.StartDt, &evt.StopDt,
		&evt.Owner, &evt.NotifyBefore, &evt.Desc); err != nil {
		return s.log.Error(ctx, "DB.DropChkScan", err)
	}
	sqlr, err := s.delete.ExecContext(ctx, map[string]any{"id": id})
	if err != nil {
		return s.log.Error(ctx, "DB.Delete", err)
	}
	n, err := sqlr.RowsAffected()
	if err == nil && n != 1 {
		return s.log.Error(ctx, "DB.Delete",
			fmt.Sprintf("unexpected rows affected count [%d]", n))
	}
	s.log.Info(ctx, "DB.Drop", "Dropped event "+evt.String())
	return nil
}

func (s *Storage) Get(ctx context.Context, dt1 time.Time, dt2 time.Time,
) ([]storage.Event, error) {
	dt1 = dt1.Truncate(time.Minute)
	dt2 = dt2.Truncate(time.Minute)
	rows, err := s.get.QueryContext(ctx, //nolint:rowserrcheck
		map[string]any{
			"dt1": dt1.Format(time.RFC3339Nano),
			"dt2": dt2.Format(time.RFC3339Nano),
		})
	if err != nil {
		return nil, s.log.Error(ctx, "DB.Get.Query", err)
	}
	defer rows.Close()
	events := []storage.Event{}
	for rows.Next() {
		var evt storage.Event
		err = rows.Scan(&evt.ID, &evt.Title, &evt.StartDt, &evt.StopDt,
			&evt.Owner, &evt.NotifyBefore, &evt.Desc)
		if err != nil {
			return nil, s.log.Error(ctx, "DB.Get.Scan", err)
		}
		events = append(events, evt)
		s.log.Debug(ctx, "DB.Get.Row", evt.String(), 5)
	}
	s.log.Debug(ctx, "DB.Get", fmt.Sprintf("Query for datetime between [%s] and [%s] = %d events selected",
		dt1.Format(time.RFC3339Nano), dt2.Format(time.RFC3339Nano), len(events)), 1)
	return events, nil
}
