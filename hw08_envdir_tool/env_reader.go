package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrInvalidCharacterInFileName = "invalid character in the file name"

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	environment := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir err: %w", err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "=") {
			return nil, errors.New(ErrInvalidCharacterInFileName)
		}

		err = readFile(dir, file.Name(), environment)
		if err != nil {
			return nil, err
		}
	}

	return environment, nil
}

func readFile(dir string, fileName string, environment Environment) (err error) {
	environment[fileName] = EnvValue{}

	file, err := os.OpenFile(path.Join(dir, fileName), os.O_RDONLY, 0o666)
	if err != nil {
		return fmt.Errorf("open file err: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			if err != nil {
				err = fmt.Errorf("%w; %w", err, closeErr)
				return
			}

			err = closeErr
		}
	}()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		environment[fileName] = EnvValue{
			Value: strings.TrimRight(string(bytes.ReplaceAll(scanner.Bytes(), []byte{0x00}, []byte("\n"))), "\t &nbsp"),
		}
	} else if err = scanner.Err(); err != nil {
		return
	} else {
		environment[fileName] = EnvValue{
			NeedRemove: true,
		}
	}

	return nil
}
