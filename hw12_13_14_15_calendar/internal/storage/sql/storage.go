package sqlstorage

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/config"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/helper"
	model "github.com/fixme_my_friend/hw12_13_14_15_calendar/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type SQLStorage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config) (*SQLStorage, error) {
	connStr := helper.ConCat(
		cfg.GetString("db.driver"), "://", cfg.GetString("db.dsn.user"), ":",
		cfg.GetString("db.dsn.password"), "@", cfg.GetString("db.dsn.host"), ":",
		cfg.GetString("db.dsn.port"), "/", cfg.GetString("db.dsn.name"), "?sslmode=",
		cfg.GetString("db.dsn.sslmode"), "&search_path=", cfg.GetString("db.dsn.search_path"),
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

	return &SQLStorage{
		pool: pool,
	}, nil
}

func (s *SQLStorage) Connect(ctx context.Context) error {
	if err := s.pool.Ping(ctx); err != nil {
		return errors.Wrap(err, "cannot connect")
	}

	return nil
}

func (s *SQLStorage) Close(ctx context.Context) {
	s.pool.Close()
}

func (s *SQLStorage) Create(ctx context.Context, event model.Event) (uint, error) {
	var id uint

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "cannot open connect")
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "cannot begin")
	}
	defer func(ctx context.Context, tx pgx.Tx) {
		if err := tx.Rollback(ctx); err != nil {
			log.Fatalln(err)
		}
	}(ctx, tx)

	params := []interface{}{
		event.Title,
		event.StartDT,
		event.EndDT,
		event.Description,
		event.UserID,
		event.NotifDT,
	}

	sqlstr := helper.ConCat(
		"INSERT INTO public.event ",
		"(title, start_dt, end_dt, description, user_id, notif_dt) ",
		"VALUES ($1, $2, $3, $4, $5, $6) ",
		"RETURNING id;",
	)

	row := conn.QueryRow(ctx, sqlstr, params)
	if err := row.Scan(&id); err != nil {
		return 0, errors.Wrap(err, "cannot insert")
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, errors.Wrap(err, "cannot commit")
	}

	return id, nil
}

func (s *SQLStorage) Update(ctx context.Context, event model.Event) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot open connect")
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot begin")
	}
	defer func(ctx context.Context, tx pgx.Tx) {
		if err := tx.Rollback(ctx); err != nil {
			log.Fatalln(err)
		}
	}(ctx, tx)

	params := []interface{}{
		event.Title,
		event.StartDT,
		event.EndDT,
		event.Description,
		event.UserID,
		event.NotifDT,
		event.ID,
	}

	sqlstr := helper.ConCat(
		"UPDATE public.event SET ",
		"title=$1, ",
		"start_dt=$2, ",
		"end_dt=$3, ",
		"description=$4, ",
		"user_id=$5, ",
		"notif_dt=$6 ",
		"WHERE id = $7;",
	)

	_, err = tx.Exec(ctx, sqlstr, params)
	if err != nil {
		return errors.Wrap(err, "cammot exec")
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "cannot commit")
	}

	return nil
}

func (s *SQLStorage) Delete(ctx context.Context, eventID uint) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot open connect")
	}
	defer conn.Release()

	sqlstr := "DELETE FROM public.event WHERE id = $1;"

	_, err = conn.Exec(ctx, sqlstr, eventID)
	if err != nil {
		return errors.Wrap(err, "cannot exec")
	}

	return nil
}

func (s *SQLStorage) GetEventByID(ctx context.Context, eventID uint) (model.Event, error) {
	var result model.Event

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return model.Event{}, errors.Wrap(err, "cannot open connect")
	}
	defer conn.Release()

	sqlstr := helper.ConCat(
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
		return model.Event{}, errors.Wrap(err, "cannot scan")
	}

	return result, nil
}

func (s *SQLStorage) GetEventsByParams(ctx context.Context, args map[string]interface{}) ([]model.Event, error) {
	var nP int
	var params []interface{}
	var result []model.Event

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open connect")
	}
	defer conn.Release()

	sqlstr := helper.ConCat(
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
		sqlstr = helper.ConCat(sqlstr, " and id = $", strconv.Itoa(nP))
		params = append(params, id)
	}

	if ids, ok := args["ids"].([]int); ok {
		var sids string
		nP++
		sqlstr = helper.ConCat(sqlstr, " and id IN ($", strconv.Itoa(nP), ")")
		for _, id := range ids {
			sids = helper.ConCat(sids, strconv.Itoa(id), ",")
		}
		params = append(params, sqlstr)
	}

	if title, ok := args["title"].(string); ok {
		nP++
		sqlstr = helper.ConCat(sqlstr, " and title = $", strconv.Itoa(nP))
		params = append(params, title)
	}

	if userID, ok := args["user_id"].(int); ok {
		nP++
		sqlstr = helper.ConCat(sqlstr, " and user_id = $", strconv.Itoa(nP))
		params = append(params, userID)
	}

	if startTime, ok := args["start_dt"].(time.Time); ok {
		nP++
		sqlstr = helper.ConCat(sqlstr, " and start_dt = $", strconv.Itoa(nP))
		params = append(params, startTime)
	}

	if endTime, ok := args["end_dt"].(time.Time); ok {
		nP++
		sqlstr = helper.ConCat(sqlstr, " and start_dt = $", strconv.Itoa(nP))
		params = append(params, endTime)
	}

	sqlstr = helper.ConCat(sqlstr, ";")

	rows, err := conn.Query(ctx, sqlstr, params)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query")
	}
	defer rows.Close()

	var e model.Event
	for rows.Next() {
		e = model.Event{}
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.StartDT,
			&e.EndDT,
			&e.Description,
			&e.UserID,
			&e.NotifDT,
		); err != nil {
			return nil, errors.Wrap(err, "cannot scan")
		}

		result = append(result, e)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "cannot rows")
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

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open connect")
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, sql, params)
	if err != nil {
		return nil, errors.Wrap(err, "cannot select")
	}
	defer rows.Close()

	for rows.Next() {
		event := model.Event{}

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.StartDT,
			&event.EndDT,
			&event.NotifDT,
			&event.UserID,
			&event.Description,
		)
		if err != nil {
			return nil, errors.Wrap(err, "cannot scan")
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(err, "cannot rows")
	}

	return events, nil
}
