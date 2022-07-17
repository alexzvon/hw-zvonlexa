package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	logger, err := New("logger_test.txt")

	require.Nil(t, err)
	require.NotNil(t, logger)

	logger.Debug("Test Debug Message")
	logger.Info("Test Info Message")
	logger.Warn("Test Warn Message")
	logger.Error("Test Error Message")

	logger.Close()
}
