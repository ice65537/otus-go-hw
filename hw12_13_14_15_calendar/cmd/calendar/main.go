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
	memstore "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	const (
		opStart = "Server.Start"
		opStop  = "Server.Stop"
	)
	var storage app.Storage // err     error

	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := GetConfig()

	switch config.Storage.Type {
	case "memory":
		storage = memstore.New()
	default:
		panic(fmt.Errorf("storage type [%s] unknown", config.Storage.Type))
	}

	app := app.New("Calendar.Listener", config.Logger.Level, config.Logger.Depth, storage)
	log := app.Logger()

	server := internalhttp.NewServer(app)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error(opStop, "failed to stop http server: "+err.Error())
		} else {
			log.Info(opStop, "Server has been stopped")
		}
	}()

	if err := server.Start(ctx); err != nil {
		log.Error(opStart, "failed to start http server: "+err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	log.Info(opStart, "Server is running...")
}
