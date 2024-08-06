package main

import (
	"fmt"
	"log"
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
	fmt.Println(address)
	fmt.Println(timeout)
}
