package main

import (
	"bufio"
	"errors"
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
	// возвращаем структуру(объект) которая имплементирует интерфейс
	// т.к. все методы интерфейса имеют реализацию
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
	if c.conn != nil {
		c.closed = true
		return c.conn.Close()
	}
	return errors.New("not opened connection")
}

func (c *client) Send() error {
	if c.conn != nil {
		// создать буферезированный reader для stdin - как самый эффективный способ чтения
		r := bufio.NewReader(c.in)
		for {
			// читать до переноса строки
			str, err := r.ReadString('\n')
			if errors.Is(err, io.EOF) {
				log.Println("...EOF")
				return nil
			}
			if err != nil {
				return c.formatSendError(err)
			}

			// записать данные в соединение
			_, err = c.conn.Write([]byte(str))
			if err != nil {
				return c.formatSendError(err)
			}
		}
	}
	return errors.New("not opened connection")
}

func (c *client) Receive() error {
	if c.conn != nil {
		// создать буферезированный reader для соединения
		r := bufio.NewReader(c.conn)
		for {
			str, err := r.ReadString('\n')
			if errors.Is(err, io.EOF) {
				log.Println("...Connection was closed")
				return nil
			}
			if err != nil {
				return c.formatReceiveError(err)
			}

			// записать данные в stdout
			_, err = c.out.Write([]byte(str))
			if err != nil {
				return c.formatReceiveError(err)
			}
		}
	}
	return errors.New("not opened connection")
}

func (c *client) formatSendError(err error) error {
	if c.closed {
		return nil
	}
	return fmt.Errorf("cannot send: %w", err)
}

func (c *client) formatReceiveError(err error) error {
	if c.closed {
		return nil
	}
	return fmt.Errorf("cannot receive: %w", err)
}
