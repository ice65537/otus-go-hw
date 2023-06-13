package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	log *logger.Logger
	srv *http.Server
}

type SrvHandler struct {
	log *logger.Logger
}

type Application interface {
	Logger() *logger.Logger
}

func NewServer(app Application, host string, port, timeout int) *Server {
	sh := SrvHandler{log: app.Logger()}
	srv := http.Server{
		Addr:         fmt.Sprintf(host+":%d", port),
		Handler:      midWare(app.Logger(),&sh),
		ReadTimeout:  time.Duration(timeout) * time.Second,
		WriteTimeout: time.Duration(timeout) * time.Second,
	}
	return &Server{log: app.Logger(), srv: &srv}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.Stop(ctx)
	}()
	s.log.Info("Server.Starting","Starting at address^ "+s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return s.log.Error("Server.Listen", fmt.Sprintf("%v", err))
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return s.log.Error("Server.Stop", fmt.Sprintf("%v", err))
	}
	return nil
}

func (sh *SrvHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!"))
}
