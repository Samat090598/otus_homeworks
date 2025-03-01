package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		dir, _ := os.Getwd()
		dir = path.Join(dir, "testdata/echo.sh")

		cmd := []string{"/bin/bash", dir, "arg1=1", "arg2=2"}

		result := RunCmd(cmd, ENV)

		require.Equal(t, 0, result)
	})

	t.Run("error case", func(t *testing.T) {
		dir, _ := os.Getwd()
		dir = path.Join(dir, "testdata/echo")

		cmd := []string{"/bin/bash", dir, "arg1=1", "arg2=2"}

		result := RunCmd(cmd, ENV)

		require.NotEqual(t, 0, result)
	})
}
