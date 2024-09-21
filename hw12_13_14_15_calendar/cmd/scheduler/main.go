package main

import (
	"context"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/app"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	internalrmqproducer "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/mq/producer/rmq"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Short: "Scheduler application",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()
		config, err := NewConfig(cfgFile)
		if err != nil {
			log.Println("Error create config: " + err.Error())
			return
		}
		logg := logger.New("info", os.Stdout)
		if err != nil {
			log.Println("Error create app logger: " + err.Error())
			return
		}
		storage := app.GetEventUseCase(config.Database)
		err = storage.Connect(ctx)
		if err != nil {
			logg.Error("Error create db connection: " + err.Error())
			return
		}
		produce := internalrmqproducer.New(config.Producer, logg)
		appl := scheduler.New(logg, produce, storage, config.Duration)
		if err = produce.Connect(ctx); err != nil {
			logg.Error("Error create mq connection: " + err.Error())
		}
		go func() {
			if err := appl.Run(ctx); err != nil {
				logg.Error("failed to consume mq: " + err.Error())
				cancel()
			}
		}()
		defer produce.Close(ctx)
		defer storage.Close(ctx)
		<-ctx.Done()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./configs/scheduler_config.toml", "Configuration file path")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
