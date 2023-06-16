package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"
	internalmem "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	internaldb "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	var storage app.Storage

	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg := GetConfig()

	app := app.New("Calendar.Listener", cfg.Logger.Level, cfg.Logger.Depth, storage)
	log := app.Logger()

	switch cfg.Storage.Type {
	case "memory":
		storage = internalmem.New()
	case "postgres":
		storage = internaldb.New(cfg.Storage.Postgre.Host,
			cfg.Storage.Postgre.Port,
			cfg.Storage.Postgre.Dbname,
			cfg.Storage.Postgre.Username,
			cfg.Storage.Postgre.Password,
		)
	default:
		panic(fmt.Errorf("storage type [%s] unknown", cfg.Storage.Type))
	}

	server := internalhttp.NewServer(app, cfg.Server.Host, cfg.Server.Port, cfg.Server.Timeout)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error(ctx, "Server.Stop", "failed to stop http server: "+err.Error())
		} else {
			log.Info(ctx, "Server.Stop", "Server has been stopped")
		}
	}()

	if err := server.Start(ctx); err != nil {
		log.Error(ctx, "Server.Start", "failed to start http server: "+err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
