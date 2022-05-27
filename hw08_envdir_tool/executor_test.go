package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Success case simple", func(t *testing.T) {
		cmd := []string{"pwd", "-L"}
		exitCode := RunCmd(cmd, Environment{})

		require.Equal(t, exitCodeOk, exitCode)
	})

	t.Run("Success case with filled dir env", func(t *testing.T) {
		cmd := []string{"pwd", "-L"}
		exitCode := RunCmd(cmd, Environment{
			"TEST_SET": EnvValue{
				Value: "TEST_SET",
			},
		})

		require.Equal(t, exitCodeOk, exitCode)
		require.Contains(t, os.Environ(), "TEST_SET=TEST_SET")
	})

	t.Run("Exec err case", func(t *testing.T) {
		cmd := []string{"pwd", "-RRR"}
		exitCode := RunCmd(cmd, Environment{
			"TEST_QWE": EnvValue{
				Value: "TEST_QWE",
			},
		})

		require.Equal(t, 1, exitCode)
	})
}
