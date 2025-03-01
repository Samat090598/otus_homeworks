package main

import (
	"errors"
	"log/slog"
	"os"
)

var ErrNotEnoughArguments = errors.New("not enough arguments")

func main() {
	args := os.Args

	if len(args) < 4 {
		slog.Error(ErrNotEnoughArguments.Error())
		os.Exit(1)
	}

	environment, err := ReadDir(os.Args[1])
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	_ = RunCmd(args[2:], environment)
}
