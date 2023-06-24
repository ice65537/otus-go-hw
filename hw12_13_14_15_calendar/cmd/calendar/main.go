package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	internalapp "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	internalgrpc "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"
	internalmem "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	internaldb "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/jackc/pgx/stdlib"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	var storage internalapp.Storage

	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg := GetConfig()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	connStr := ""
	switch cfg.Storage.Type {
	case "memory":
		storage = internalmem.New()
	case "postgres":
		storage = internaldb.New()
		connStr = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
			cfg.Storage.Postgre.Host,
			cfg.Storage.Postgre.Port,
			cfg.Storage.Postgre.Dbname,
			cfg.Storage.Postgre.Username,
			cfg.Storage.Postgre.Password,
		)
	default:
		panic(fmt.Errorf("storage type [%s] unknown", cfg.Storage.Type))
	}

	app := internalapp.New("Calendar.Keeper", cfg.Logger.Level, cfg.Logger.Depth, storage, cancel)
	if err := app.Init(ctx, connStr); err != nil {
		cancel()
		os.Exit(1) //nolint:gocritic,nolintlint
	}

	servHTTP := internalhttp.NewServer(app, cfg.ServerHTTP.Host, cfg.ServerHTTP.Port, cfg.ServerHTTP.Timeout)
	servGRPC := internalgrpc.NewServer(app, cfg.ServerGRPC.Host, cfg.ServerGRPC.Port, cfg.ServerGRPC.Timeout)

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		_ = servHTTP.Stop(ctx)
	}()

	go func() {
		<-ctx.Done()
		// ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		servGRPC.Stop()
	}()

	go func() {
		if err := servHTTP.Start(ctx); err != nil {
			cancel()
			os.Exit(1) //nolint:gocritic,nolintlint
		}
	}()
	go func() {
		if err := servGRPC.Start(ctx); err != nil {
			cancel()
			os.Exit(1) //nolint:gocritic,nolintlint
		}
	}()
	<-ctx.Done()
}
