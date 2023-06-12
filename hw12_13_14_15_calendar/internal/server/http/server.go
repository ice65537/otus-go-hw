package internalhttp

import (
	"context"
	"net/http"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	log *logger.Logger
}

type Application interface {
	Logger() *logger.Logger
}

func NewServer(app Application) *Server {
	return &Server{log: app.Logger()}
}

func (s *Server) Start(ctx context.Context) error {
	http.HandleFunc("/",)

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
