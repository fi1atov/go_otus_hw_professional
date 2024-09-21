package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	internalrmqconsumer "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/mq/consumer/rmq"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/sender"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Short: "Sender application",
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
		fmt.Println("OKOKOKO")
		consumer := internalrmqconsumer.New(config.Consumer, logg)
		app := sender.New(logg, consumer)
		if err := consumer.Connect(ctx); err != nil {
			logg.Error("Error create mq connection: " + err.Error())
		}
		go func() {
			if err := app.Consume(ctx); err != nil {
				logg.Error("failed to consume mq: " + err.Error())
				cancel()
			}
		}()
		defer consumer.Close(ctx)
		<-ctx.Done()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./configs/config_sender.toml", "Configuration file path")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
