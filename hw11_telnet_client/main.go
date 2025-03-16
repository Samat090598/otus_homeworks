package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")

	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		slog.Error("not enough arguments")
		return
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		slog.Error(fmt.Errorf("connect err: %w", err).Error())
		return
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error(fmt.Errorf("close err: %w", err).Error())
			return
		}
	}()

	errCh := make(chan error)
	go func() {
		errCh <- client.Send()
	}()
	go func() {
		errCh <- client.Receive()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		slog.Error(err.Error())
	}
}
