package main

import (
	"context"
	"flag"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage/initstoragesql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/app"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/server/http"
)

// go run cmd/calendar/*.go --config=configs/config.toml

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	mainCtx, cancel := context.WithCancel(context.Background())
	go watchSignals(cancel)

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg := logger.New(config.Logger.Level, os.Stdout, config.Logger.File)

	logg.Info("start calendar")

	db, err := initstoragesql.New(mainCtx, config.Database.Inmem, config.Database.Connect)
	if err != nil {
		logg.Error("err")
	}
	calendar := app.New(logg, db)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func watchSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}
