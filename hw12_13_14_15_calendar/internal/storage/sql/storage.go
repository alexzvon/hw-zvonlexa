package sqlstorage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/config"
	"github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/myutils"
	model "github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type SQLStorage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config) (*SQLStorage, error) {
	connStr := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
		cfg.GetString("db.driver"), cfg.GetString("db.dsn.user"),
		cfg.GetString("db.dsn.password"), cfg.GetString("db.dsn.host"),
		cfg.GetString("db.dsn.port"), cfg.GetString("db.dsn.name"),
		cfg.GetString("db.dsn.sslmode"), cfg.GetString("db.dsn.search_path"),
	)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.GetInt("db.postgres.maxconns"))
	poolConfig.MinConns = int32(cfg.GetInt("db.postgres.minconns"))

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "cannot connect")
	}

	return &SQLStorage{
		pool: pool,
	}, nil
}

func (s *SQLStorage) Close(ctx context.Context) {
	s.pool.Close()
}

func (s *SQLStorage) Create(ctx context.Context, event model.Event) (uint, error) {
	var id uint

	f := func(conn *pgxpool.Conn) error {
		params := []interface{}{
			event.Title,
			event.StartDT,
			event.EndDT,
			event.Description,
			event.UserID,
			event.NotifDT,
		}

		sqlstr := myutils.ConCat(
			"INSERT INTO public.event ",
			"(title, start_dt, end_dt, description, user_id, notif_dt) ",
			"VALUES ($1, $2, $3, $4, $5, $6) ",
			"RETURNING id;",
		)

		row := conn.QueryRow(ctx, sqlstr, params)
		if err := row.Scan(&id); err != nil {
			return errors.Wrap(err, "cannot insert")
		}

		return nil
	}

	if err := s.pool.AcquireFunc(ctx, f); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SQLStorage) Update(ctx context.Context, event model.Event) error {
	f := func(tx pgx.Tx) error {
		params := []interface{}{
			event.Title,
			event.StartDT,
			event.EndDT,
			event.Description,
			event.UserID,
			event.NotifDT,
			event.ID,
		}

		sqlstr := myutils.ConCat(
			"UPDATE public.event SET ",
			"title=$1, ",
			"start_dt=$2, ",
			"end_dt=$3, ",
			"description=$4, ",
			"user_id=$5, ",
			"notif_dt=$6 ",
			"WHERE id = $7;",
		)

		if _, err := tx.Exec(ctx, sqlstr, params); err != nil {
			return errors.Wrap(err, "cannot update")
		}

		return nil
	}

	if err := s.pool.BeginFunc(ctx, f); err != nil {
		return err
	}

	return nil
}

func (s *SQLStorage) Delete(ctx context.Context, eventID uint) error {
	f := func(conn *pgxpool.Conn) error {
		sqlstr := "DELETE FROM public.event WHERE id = $1;"

		if _, err := conn.Exec(ctx, sqlstr, eventID); err != nil {
			return errors.Wrap(err, "cannot delete")
		}

		return nil
	}

	if err := s.pool.AcquireFunc(ctx, f); err != nil {
		return err
	}

	return nil
}

func (s *SQLStorage) GetEventByID(ctx context.Context, eventID uint) (model.Event, error) {
	var result model.Event

	f := func(conn *pgxpool.Conn) error {
		sqlstr := myutils.ConCat(
			"SELECT ",
			"id, ",
			"title, ",
			"start_dt, ",
			"end_dt, ",
			"description, ",
			"user_id, ",
			"notif_dt, ",
			"FROM public.event WHERE id = $1",
		)

		row := conn.QueryRow(ctx, sqlstr, eventID)
		if err := row.Scan(&result); err != nil {
			return errors.Wrap(err, "cannot GetEventByID")
		}

		return nil
	}

	if err := s.pool.AcquireFunc(ctx, f); err != nil {
		return model.Event{}, err
	}

	return result, nil
}

func (s *SQLStorage) GetEventsByParams(ctx context.Context, args map[string]interface{}) ([]model.Event, error) {
	var nP int
	var params []interface{}
	var result []model.Event

	event := model.Event{}
	scans := []interface{}{
		&event.ID,
		&event.Title,
		&event.StartDT,
		&event.EndDT,
		&event.Description,
		&event.UserID,
		&event.NotifDT,
	}

	sqlstr := myutils.ConCat(
		"SELECT ",
		"id, ",
		"title, ",
		"start_dt, ",
		"end_dt, ",
		"description, ",
		"user_id, ",
		"notif_dt ",
		"FROM public.event WHERE 1 ",
	)

	if id, ok := args["id"].(int); ok {
		nP++
		sqlstr = myutils.ConCat(sqlstr, " and id = $", strconv.Itoa(nP))
		params = append(params, id)
	}

	if ids, ok := args["ids"].([]int); ok {
		var sids string
		nP++
		sqlstr = myutils.ConCat(sqlstr, " and id IN ($", strconv.Itoa(nP), ")")
		for _, id := range ids {
			sids = myutils.ConCat(sids, strconv.Itoa(id), ",")
		}
		params = append(params, sqlstr)
	}

	if title, ok := args["title"].(string); ok {
		nP++
		sqlstr = myutils.ConCat(sqlstr, " and title = $", strconv.Itoa(nP))
		params = append(params, title)
	}

	if userID, ok := args["user_id"].(int); ok {
		nP++
		sqlstr = myutils.ConCat(sqlstr, " and user_id = $", strconv.Itoa(nP))
		params = append(params, userID)
	}

	if startTime, ok := args["start_dt"].(time.Time); ok {
		nP++
		sqlstr = myutils.ConCat(sqlstr, " and start_dt = $", strconv.Itoa(nP))
		params = append(params, startTime)
	}

	if endTime, ok := args["end_dt"].(time.Time); ok {
		nP++
		sqlstr = myutils.ConCat(sqlstr, " and start_dt = $", strconv.Itoa(nP))
		params = append(params, endTime)
	}

	sqlstr = myutils.ConCat(sqlstr, ";")

	f := func(pgx.QueryFuncRow) error {
		result = append(result, event)
		return nil
	}

	if _, err := s.pool.QueryFunc(ctx, sqlstr, params, scans, f); err != nil {
		return nil, errors.Wrap(err, "cannot GetEventsByParams")
	}

	return result, nil
}

func (s *SQLStorage) ListEventsToDay(ctx context.Context, dt time.Time) ([]model.Event, error) {
	sql := "SELECT * FROM public.event WHERE start_dt=$1;"

	params := make([]interface{}, 0)
	params = append(params, dt)

	events, err := s.rowsSelect(ctx, sql, params)
	if err != nil {
		return nil, errors.Wrap(err, "cannot rows select")
	}

	return events, nil
}

func (s *SQLStorage) ListEventsToWeek(ctx context.Context, dt time.Time) ([]model.Event, error) {
	sql := "SELECT * FROM public.event WHERE start_dt<$1 AND start_dt>$2;"

	params := make([]interface{}, 0)
	params = append(params, dt)
	params = append(params, dt.Add(7))

	events, err := s.rowsSelect(ctx, sql, params)
	if err != nil {
		return nil, errors.Wrap(err, "cannot rows select")
	}

	return events, nil
}

func (s *SQLStorage) ListEventsToMonth(ctx context.Context, dt time.Time) ([]model.Event, error) {
	sql := "SELECT * FROM public.event WHERE start_dt<$1 AND start_dt>$2;"

	params := make([]interface{}, 0)
	params = append(params, dt)
	params = append(params, dt.Add(30))

	events, err := s.rowsSelect(ctx, sql, params)
	if err != nil {
		return nil, errors.Wrap(err, "cannot rows select")
	}

	return events, nil
}

func (s *SQLStorage) rowsSelect(ctx context.Context, sql string, params []interface{}) ([]model.Event, error) {
	var events []model.Event

	event := model.Event{}
	scans := []interface{}{
		&event.ID,
		&event.Title,
		&event.StartDT,
		&event.EndDT,
		&event.Description,
		&event.UserID,
		&event.NotifDT,
	}

	f := func(pgx.QueryFuncRow) error {
		events = append(events, event)
		return nil
	}

	if _, err := s.pool.QueryFunc(ctx, sql, params, scans, f); err != nil {
		return nil, errors.Wrap(err, "cannot rows select")
	}

	return events, nil
}
