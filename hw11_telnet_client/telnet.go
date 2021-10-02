package main

import (
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
		return err
	}
	c.conn = conn
	return nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Send() error {
	buff := make([]byte, 1024)
	n, err := c.in.Read(buff)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(buff[0:n])
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Receive() error {
	buff := make([]byte, 1024)
	n, err := c.conn.Read(buff)
	if err != nil {
		return err
	}
	_, err = c.out.Write(buff[0:n])
	if err != nil {
		return err
	}
	return nil
}
