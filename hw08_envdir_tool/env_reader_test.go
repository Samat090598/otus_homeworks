package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

var ENV = Environment{
	"EMPTY": EnvValue{
		Value: "",
	},
	"FOO": EnvValue{
		Value: `   foo
with new line`,
	},
	"HELLO": EnvValue{
		Value: `"hello"`,
	},
	"UNSET": EnvValue{
		NeedRemove: true,
	},
	"BAR": EnvValue{
		Value: `bar`,
	},
}

func TestReadDir(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		dir, _ := os.Getwd()

		dir = path.Join(dir, "testdata/env")

		environment := Environment{
			"EMPTY": EnvValue{
				Value: "",
			},
			"FOO": EnvValue{
				Value: `   foo
with new line`,
			},
			"HELLO": EnvValue{
				Value: `"hello"`,
			},
			"UNSET": EnvValue{
				NeedRemove: true,
			},
			"BAR": EnvValue{
				Value: `bar`,
			},
		}

		result, err := ReadDir(dir)

		require.NoError(t, err)
		require.Equal(t, environment, result)
	})

	t.Run("incorrect dir", func(t *testing.T) {
		dir, _ := os.Getwd()
		dir = path.Join(dir, "test")

		_, err := ReadDir(dir)
		require.Error(t, err)
	})
}
