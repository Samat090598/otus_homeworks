package main

import (
	"flag"
	"log/slog"
	"net"
	"os"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")

	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		slog.Error("not enough arguments")
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		slog.Error("connect err:", err)
		os.Exit(1)
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("close err:", err)
			os.Exit(1)
		}
	}()

	go func() {
		if err := client.Receive(); err != nil {
			slog.Error("receive err:", err)
			os.Exit(1)
		}
	}()

	if err := client.Send(); err != nil {
		slog.Error("send err:", err)
		os.Exit(1)
	}
}
