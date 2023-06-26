package internalgrpc

import (
	"context"
	"fmt"
	"net"

	eventsrv "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/api"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc"
)

type Server struct {
	eventsrv.UnimplementedEventsServer
	app     server.Application
	log     *logger.Logger
	grpcSrv *grpc.Server
	addr    string
	timeout int
}

func NewServer(app server.Application, host string, port, timeout int) *Server {
	appSrv := &Server{
		log:     app.Logger(),
		app:     app,
		addr:    fmt.Sprintf("%s:%d", host, port),
		timeout: timeout,
	}
	appSrv.grpcSrv = grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
	)
	eventsrv.RegisterEventsServer(appSrv.grpcSrv, appSrv)
	return appSrv
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.Stop()
	}()
	s.log.Info(ctx, "ServerGRPC.Start", "Starting at address "+s.addr)

	lc := net.ListenConfig{}
	lsnr, err := lc.Listen(ctx, "tcp", s.addr)
	if err != nil {
		return s.log.Fatal(ctx, "ServerGRPC.Listen", fmt.Sprintf("%v", err))
	}

	if err = s.grpcSrv.Serve(lsnr); err != nil {
		return s.log.Fatal(ctx, "ServerGRPC.Serve", fmt.Sprintf("%v", err))
	}
	return nil
}

func (s *Server) Stop() {
	s.grpcSrv.GracefulStop()
}
