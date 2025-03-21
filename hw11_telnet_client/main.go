package main

import (
	"context"
	"flag"
	"log"
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
		log.Fatal("not enough arguments")
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("connect err: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("close err: %v", err)
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
		log.Printf("error: %v", err)
	}
}
