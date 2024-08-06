package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	closed  bool
}

func (c *client) Connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("cannot connect: %w", err)
	}
	log.Printf("...Connected to %s\n", c.address)

	return nil
}

func (c *client) Close() error {
	c.closed = true
	return c.conn.Close()
}

func (c *client) Send() error {
	return nil
}

func (c *client) Receive() error {
	return nil
}
