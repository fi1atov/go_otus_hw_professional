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

	ctx, cancel := context.WithCancel(context.Background())

	go listenSignals(cancel)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	go send(client, cancel)
	go receive(client, cancel)

	<-ctx.Done()
}

func listenSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	// SIGINT посылается при нажатии Ctrl+C
	// SIGTERM посылается при использовании команды kill
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}

func send(client TelnetClient, cancel context.CancelFunc) {
	err := client.Send()
	if err != nil {
		log.Println(err)
	}
	cancel()
}

func receive(client TelnetClient, cancel context.CancelFunc) {
	err := client.Receive()
	if err != nil {
		log.Println(err)
	}
	cancel()
}
