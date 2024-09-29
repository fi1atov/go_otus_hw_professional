package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/app"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	grpcserver "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/server/http"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage/storecreator"
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

	logg := logger.New(config.Logger.Level, os.Stdout)

	logg.Info("start calendar")

	db, err := storecreator.New(mainCtx, config.App.Database)
	if err != nil {
		logg.Error("err")
	}
	calendar := app.New(logg, db)

	httpServer := internalhttp.NewServer(calendar, logg)
	go func() {
		err := httpServer.Start(config.App.HTTPServer.Host + ":" + config.App.HTTPServer.Port)
		if err != nil {
			logg.Error("Error: ", err)
			cancel()
		}
	}()

	grpcServer := grpcserver.NewServer(calendar, logg)
	go func() {
		err := grpcServer.Start(config.App.GRPCServer.Host + ":" + config.App.GRPCServer.Port)
		if err != nil {
			logg.Error("Error: ", err)
			cancel()
		}
	}()

	logg.Info("calendar is running...")

	<-mainCtx.Done()

	logg.Info("stopping calendar...")
	cancel()
	shutDown(logg, httpServer, grpcServer, db)
	logg.Info("calendar is stopped")
}

func watchSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}

func shutDown(logg logger.Logger, httpServer internalhttp.Server, grpcServer grpcserver.Server, db storage.Storage) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Stop(ctx); err != nil {
			logg.Error("Error: ", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Stop(ctx); err != nil {
			logg.Error("Error: ", err)
		}
	}()

	wg.Wait()

	if err := db.Close(ctx); err != nil {
		logg.Error("Error: ", err)
	}
}
