package internalgrpc

import (
	"context"
	"net"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/config"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/logger"
	grpcserver "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/server/grpc/gen"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
	model "github.com/fixme_my_friend/hw12_13_14_15_calendar/models"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ctxKey string

func (c ctxKey) String() string {
	return string(c)
}

type ctxVal string

func (c ctxVal) String() string {
	return string(c)
}

var (
	keyCtx   = ctxKey("name")
	valDay   = ctxVal("Day")
	valWeek  = ctxVal("Week")
	valMonth = ctxVal("Month")
)

type ServerGRPC struct {
	grpcserver.UnimplementedGRPCServerServer
	logg logger.Logger
	db   storage.Conn
	gs   *grpc.Server
}

func NewServerGRPC(cfg config.Config, logger logger.Logger, db storage.Conn) *ServerGRPC {
	s := new(ServerGRPC)
	s.logg = logger
	s.db = db
	s.gs = grpc.NewServer(unaryInterceptor())

	return s
}

func unaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		defer log.Infof(
			"%s %v",
			info.FullMethod,
			time.Since(start),
		)

		resp, err := handler(ctx, req)
		if err != nil {
			log.Errorf("method %q throws error: %v", info.FullMethod, err)
		}

		return resp, err
	})
}

func (s *ServerGRPC) UpServerGRPC(cfg config.Config, logger logger.Logger, db storage.Conn) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(cfg.GetString("server.host"), cfg.GetString("server.port")))
	if err != nil {
		return err
	}

	grpcserver.RegisterGRPCServerServer(s.gs, s)

	return s.gs.Serve(lis)
}

func (s *ServerGRPC) DownServerGRPC() {
	s.gs.GracefulStop()
}

func (s *ServerGRPC) CreateEvent(ctx context.Context, event *grpcserver.Event) (*grpcserver.IDEvent, error) {
	e := model.Event{}

	e.Title = event.GetTitle()
	e.UserID = uint(event.GetUserId())
	e.Description = event.GetDescription()

	if err := event.GetStartTime().CheckValid(); err != nil {
		return nil, err
	}
	e.StartDT = event.GetStartTime().AsTime()

	if err := event.GetEndTime().CheckValid(); err != nil {
		return nil, err
	}
	e.EndDT = event.GetEndTime().AsTime()

	if err := event.GetNotificationTime().CheckValid(); err != nil {
		return nil, err
	}
	e.NotifDT = event.GetNotificationTime().AsTime()

	id, err := s.db.Create(ctx, e)
	if err != nil {
		return nil, err
	}

	r := &grpcserver.IDEvent{Id: uint64(id)}

	return r, nil
}

func (s *ServerGRPC) UpdateEvent(ctx context.Context, event *grpcserver.Event) (*empty.Empty, error) {
	e := model.Event{}

	e.ID = uint(event.GetId())
	e.Title = event.GetTitle()
	e.UserID = uint(event.GetUserId())
	e.Description = event.GetDescription()

	if err := event.GetStartTime().CheckValid(); err != nil {
		return nil, err
	}
	e.StartDT = event.GetStartTime().AsTime()

	if err := event.GetEndTime().CheckValid(); err != nil {
		return nil, err
	}
	e.EndDT = event.GetEndTime().AsTime()

	if err := event.GetNotificationTime().CheckValid(); err != nil {
		return nil, err
	}
	e.NotifDT = event.GetNotificationTime().AsTime()

	if err := s.db.Update(ctx, e); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerGRPC) DeleteEvent(ctx context.Context, id *grpcserver.IDEvent) (*empty.Empty, error) {
	if err := s.db.Delete(ctx, uint(id.GetId())); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerGRPC) ListEventsToDay(ctx context.Context, dt *grpcserver.DateTime) (*grpcserver.Events, error) {
	ctx = context.WithValue(ctx, keyCtx, valDay)

	events, err := s.listEvents(ctx, dt)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *ServerGRPC) ListEventsToWeek(ctx context.Context, dt *grpcserver.DateTime) (*grpcserver.Events, error) {
	ctx = context.WithValue(ctx, keyCtx, valWeek)

	events, err := s.listEvents(ctx, dt)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *ServerGRPC) ListEventsToMonth(ctx context.Context, dt *grpcserver.DateTime) (*grpcserver.Events, error) {
	ctx = context.WithValue(ctx, keyCtx, valMonth)

	events, err := s.listEvents(ctx, dt)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *ServerGRPC) listEvents(ctx context.Context, dt *grpcserver.DateTime) (*grpcserver.Events, error) {
	var events []model.Event
	var err error

	switch ctx.Value(keyCtx).(string) {
	case valDay.String():
		events, err = s.db.ListEventsToDay(ctx, dt.Dt.AsTime())
	case valWeek.String():
		events, err = s.db.ListEventsToWeek(ctx, dt.Dt.AsTime())
	case valMonth.String():
		events, err = s.db.ListEventsToMonth(ctx, dt.Dt.AsTime())
	}

	if err != nil {
		return nil, err
	}

	grpcEvents := new(grpcserver.Events)

	for _, e := range events {
		event := &grpcserver.Event{
			Id:               uint64(e.ID),
			Title:            e.Title,
			StartTime:        timestamppb.New(e.StartDT),
			EndTime:          timestamppb.New(e.EndDT),
			UserId:           uint64(e.UserID),
			NotificationTime: timestamppb.New(e.NotifDT),
			Description:      e.Description,
		}

		grpcEvents.Event = append(grpcEvents.Event, event)
	}

	return grpcEvents, nil
}
