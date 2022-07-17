package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	logger, err := New("logger_test.txt")

	require.Nil(t, err)
	require.NotNil(t, logger)

	err = logger.Debug("Test Debug Message")
	require.Nil(t, err)

	err = logger.Info("Test Info Message")
	require.Nil(t, err)

	err = logger.Warn("Test Warn Message")
	require.Nil(t, err)

	err = logger.Error("Test Error Message")
	require.Nil(t, err)

	logger.Close()
}
