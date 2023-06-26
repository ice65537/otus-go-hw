package internalgrpc

import (
	"context"
	"testing"
	"time"

	eventsrv "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/api"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
	t.Run("Get hello from GRPC server", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		app := &testApplication{log: logger.New("HTTP.Test", "DEBUG", 5, cancel)}
		srv := NewServer(app, "localhost", 8888, 3)
		go func() {
			srv.Start(ctx)
		}()
		time.Sleep(3 * time.Second)

		conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure()) //nolint
		require.NoError(t, err)
		clnt := eventsrv.NewEventsClient(conn)

		md := metadata.New(nil)
		md.Append("user", "John Doe")
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		msg, err := clnt.Hello(ctx, &emptypb.Empty{})
		require.NoError(t, err)
		require.Equal(t, "Hello John Doe!", msg.Text)

		md = metadata.New(nil)
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		_, err = clnt.Hello(ctx, &emptypb.Empty{})
		x, _ := status.FromError(err)
		require.Equal(t, codes.Unauthenticated, x.Code())

		srv.Stop()
	})
}
