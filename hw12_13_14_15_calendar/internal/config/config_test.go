package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg, err := New("config_test.yaml")

	require.Nil(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, "localhost", cfg.GetString("server.host"))
	require.Equal(t, 8889, cfg.GetInt("server.port"))

	require.Equal(t, "log/log.txt", cfg.GetString("logger.path"))
	require.Equal(t, "debug", cfg.GetString("logger.level"))

	require.Equal(t, "postgres", cfg.GetString("db.postgres.driver"))
	require.Equal(t, "localhost", cfg.GetString("db.postgres.dsn.host"))
	require.Equal(t, 5401, cfg.GetInt("db.postgres.dsn.port"))
	require.Equal(t, "postgres", cfg.GetString("db.postgres.dsn.user"))
	require.Equal(t, "postgres", cfg.GetString("db.postgres.dsn.password"))
	require.Equal(t, "disable", cfg.GetString("db.postgres.sslmode"))
	require.Equal(t, 10, cfg.GetInt("db.postgres.maxconns"))
	require.Equal(t, 3, cfg.GetInt("db.postgres.minconns"))

	require.Equal(t, 20, cfg.GetInt("db.memory.maxsize"))

	require.Equal(t, "memory", cfg.GetString("repository.type"))
}
