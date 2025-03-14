package main

import (
	"errors"
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
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("dial err: %w", err)
	}
	t.conn = conn

	return nil
}

func (t *telnetClient) Close() error {
	if t.conn != nil {
		return errors.New("there is no active connection")
	}

	return t.conn.Close()
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return errors.New("there is no active connection")
	}

	_, err := io.Copy(t.conn, t.in)

	return err
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return errors.New("there is no active connection")
	}

	_, err := io.Copy(t.out, t.conn)

	return err
}
