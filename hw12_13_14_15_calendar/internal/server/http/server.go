package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/server"
)

type appGetEvents struct {
	T1 time.Time `json:"t1"`
	T2 time.Time `json:"t2"`
}

type Server struct {
	app     server.Application
	log     *logger.Logger
	httpSrv *http.Server
}

func NewServer(app server.Application, host string, port, timeout int) *Server {
	appSrv := &Server{log: app.Logger(), app: app}

	mux := http.NewServeMux()
	mux.HandleFunc("/", appSrv.hello)
	mux.HandleFunc("/hello", appSrv.hello)
	mux.HandleFunc("/event/new", appSrv.new)
	mux.HandleFunc("/event/reset", appSrv.reset)
	mux.HandleFunc("/event/drop", appSrv.drop)
	mux.HandleFunc("/event/get", appSrv.get)

	appSrv.httpSrv = &http.Server{
		Addr:         fmt.Sprintf(host+":%d", port),
		Handler:      midWareHandler(app.Logger(), mux),
		ReadTimeout:  time.Duration(timeout) * time.Second,
		WriteTimeout: time.Duration(timeout) * time.Second,
	}
	return appSrv
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		_ = s.Stop(ctx)
	}()
	s.log.Info(ctx, "Server.Starting", "Starting at address "+s.httpSrv.Addr)
	if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return s.log.Fatal(ctx, "Server.Listen", fmt.Sprintf("%v", err))
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		return s.log.Error(ctx, "Server.Stop", fmt.Sprintf("%v", err))
	}
	return nil
}
