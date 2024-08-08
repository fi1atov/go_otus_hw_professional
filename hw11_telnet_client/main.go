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
	pflag.DurationVarP(&timeout, "timeout", "t", 10*time.Second, "connection timeout")
	pflag.StringVarP(&host, "host", "h", "", "host")
	pflag.StringVarP(&port, "port", "p", "", "port")

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

	// Background - контекст - это корневой контекст который порождает все остальные контексты.
	// "Матрешка контекстов"
	// WithCancel возвращает функцию cancel которую надо вызвать чтобы прервать программу
	// Вызвать cancel() чтобы дальнейшее исполнение кода не нужно
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

	// блокируем функцию main - ждем пока не будет вызван cancel()
	// координация корректного завершения работы программы при работе
	// с несколькими горутинами
	<-ctx.Done()
}

// слушатель системных сигналов - при получении команд вырубает программу.
func listenSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	// SIGINT посылается при нажатии Ctrl+C
	// SIGTERM посылается при использовании команды kill
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}

// передача данных через клиент.
func send(client TelnetClient, cancel context.CancelFunc) {
	err := client.Send()
	if err != nil {
		log.Println(err)
	}
	cancel()
}

// получение данных через клиент.
func receive(client TelnetClient, cancel context.CancelFunc) {
	err := client.Receive()
	if err != nil {
		log.Println(err)
	}
	cancel()
}
