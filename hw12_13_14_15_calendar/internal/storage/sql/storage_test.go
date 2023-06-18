package dbstore

/*
import (
	"context"
	"testing"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("Create event, read it, update it, drop it, read nothing", func(t *testing.T) {
		store := New()
		ctx, cancel := context.WithCancel(context.Background())
		log := logger.New("test", "DEBUG", 5, cancel)
		err := store.Init(ctx, log, "host=localhost port=5432 dbname=calendar user=clndr password=clndr")
		require.NoError(t, err)

		t1 := time.Now()
		err = store.Upsert(ctx, storage.Event{
			Title:   "Test-1",
			StartDt: time.Now(),
			StopDt:  time.Now().Add(30 * time.Minute),
			Owner:   "tester",
		})
		t2 := time.Now()
		require.NoError(t, err)

		events, err2 := store.Get(ctx, t1, t2)
		require.NoError(t, err2)
		require.Equal(t, 1, len(events))
		require.Equal(t, "Test-1", events[0].Title)
		require.Equal(t, "tester", events[0].Owner)
		require.Equal(t, "", events[0].Desc)

		events[0].Desc = "Abrakadabra"
		err = store.Upsert(ctx, events[0])
		require.NoError(t, err)

		events, err = store.Get(ctx, t1, t2)
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, "Test-1", events[0].Title)
		require.Equal(t, "tester", events[0].Owner)
		require.Equal(t, "Abrakadabra", events[0].Desc)

		err = store.Drop(ctx, events[0].ID)
		require.NoError(t, err)
		events, err = store.Get(ctx, t1, t2)
		require.NoError(t, err)
		require.Equal(t, 0, len(events))
	})
}
*/
