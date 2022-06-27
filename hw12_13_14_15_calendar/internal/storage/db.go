package storage

import (
	"context"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/config"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/model"
	memorystorage "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pkg/errors"
)

const (
	RepTypeMemory   = "memory"
	RepTypePostgres = "postgres"
)

type Conn interface {
	Create(ctx context.Context, event model.Event) (uint, error)
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, eventID uint) error
	GetEventByID(ctx context.Context, eventID uint) (model.Event, error)
	GetEventsByParams(ctx context.Context, args map[string]interface{}) ([]model.Event, error)
	Close(ctx context.Context)
}

func Connect(cfg config.Config) (Conn, error) {
	var db interface{}
	var err error

	switch cfg.GetString("repository.type") {
	case RepTypeMemory:
		db, err = memorystorage.New(cfg.GetInt("db.memory.maxsize"))
		if err != nil {
			return nil, errors.New("cannot connect memory repository")
		}
	case RepTypePostgres:
		db, err = sqlstorage.New(context.Background(), cfg)
		if err != nil {
			return nil, errors.New("cannot connect postgres repository")
		}
	default:
		return nil, errors.New("cannot create repository")
	}

	return db.(Conn), nil
}
