package main

import (
	"fmt"
	"io"
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
		in:      in,
		out:     out,
		timeout: timeout,
	}
}

type client struct {
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	address string
	timeout time.Duration
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	c.conn = conn
	return nil
}

func (c *client) Close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("close: %w", err)
		}
	}
	return nil
}

func (c *client) Send() error {
	_, err := io.Copy(c.conn, c.in)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (c *client) Receive() error {
	_, err := io.Copy(c.out, c.conn)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}
	return nil
}
