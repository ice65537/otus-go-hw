package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

type Application interface {
	Logger() *logger.Logger
}

type Server struct {
	log     *logger.Logger
	httpSrv *http.Server
}

func NewServer(app Application, host string, port, timeout int) *Server {
	appSrv := &Server{log: app.Logger()}

	mux := http.NewServeMux()
	mux.HandleFunc("/", appSrv.hello)
	mux.HandleFunc("/hello", appSrv.hello)
	mux.HandleFunc("/bye", appSrv.hello)

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
		s.Stop(ctx)
	}()
	s.log.Info(ctx, "Server.Starting", "Starting at address "+s.httpSrv.Addr)
	if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return s.log.Error(ctx, "Server.Listen", fmt.Sprintf("%v", err))
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		return s.log.Error(ctx, "Server.Stop", fmt.Sprintf("%v", err))
	}
	return nil
}

func (s *Server) hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello %s!", getReqSession(r).User)))
}
