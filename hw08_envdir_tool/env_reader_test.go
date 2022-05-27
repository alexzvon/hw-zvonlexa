package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("test is valid", func(t *testing.T) {
		expectedEnv := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		env, err := ReadDir("testdata/env")

		require.NoError(t, err)
		require.Equal(t, expectedEnv, env)
	})

	t.Run("test when empty dir", func(t *testing.T) {
		testDir, err := os.MkdirTemp("", "empty_dir")
		if err != nil {
			t.Fatal("can't create temp dir: ", err)
		}

		res, err := ReadDir(testDir)

		require.NoError(t, err)
		require.Len(t, res, 0)
	})

	t.Run("test when non existent dir", func(t *testing.T) {
		res, err := ReadDir("testdata/312321321/")

		require.Error(t, err)
		require.Equal(t, true, os.IsNotExist(err))
		require.Len(t, res, 0)
	})
}
