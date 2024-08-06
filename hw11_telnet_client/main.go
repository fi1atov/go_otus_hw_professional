package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

func getParams() (address string, timeout time.Duration) {
	var host, port string
	pflag.StringVarP(&host, "host", "h", "", "host")
	pflag.StringVarP(&port, "port", "p", "", "port")
	pflag.DurationVarP(&timeout, "timeout", "t", 10*time.Second, "connection timeout")

	pflag.Parse()
	if host == "" || port == "" {
		log.Fatal("Please define address and port")
	}
	address = host + ":" + port
	return
}

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	address, timeout := getParams()

	_, cancel := context.WithCancel(context.Background())

	go listenSignals(cancel)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
}

func listenSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	// SIGINT посылается при нажатии Ctrl+C
	// SIGTERM посылается при использовании команды kill
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}
