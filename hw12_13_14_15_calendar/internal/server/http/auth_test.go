package internalhttp

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/require"
)

type testApplication struct {
	log *logger.Logger
	evt storage.Event
}

func (ta *testApplication) Logger() *logger.Logger {
	return ta.log
}

func (ta *testApplication) Upsert(ctx context.Context, evt storage.Event) error {
	ta.evt = evt
	return nil
}

func (ta *testApplication) Drop(ctx context.Context, evtId string) error {
	ta.evt = storage.Event{}
	return nil
}

func (ta *testApplication) Get(ctx context.Context, t1 time.Time, t2 time.Time) ([]storage.Event, error) {
	evts := make([]storage.Event, 2)
	evts[0] = ta.evt
	evts[1].StartDt = t1
	evts[2].StopDt = t2
	return evts, nil
}

func TestHello(t *testing.T) {
	t.Run("Get hello from HTTP server", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		app := &testApplication{log: logger.New("HTTP.Test", "DEBUG", 5, cancel)}
		srv := NewServer(app, "localhost", 7777, 3)
		go func() {
			srv.Start(ctx)
		}()
		time.Sleep(3 * time.Second)
		clnt := &http.Client{}

		req, _ := http.NewRequest("GET", "http://localhost:7777", nil)
		req.Header.Add("User", "John Doe")
		resp, err := clnt.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		respB, _ := io.ReadAll(resp.Body)
		require.Equal(t, "Hello John Doe!", string(respB))

		req, _ = http.NewRequest("GET", "http://localhost:7777", nil)
		resp, err = clnt.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, resp.StatusCode)

		err = srv.Stop(ctx)
		require.NoError(t, err)
	})
}
