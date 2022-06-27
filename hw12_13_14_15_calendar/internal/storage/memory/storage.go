package memorystorage

import (
	"context"
	"sync"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/model"
	"github.com/pkg/errors"
)

type MemStorage struct {
	mu      sync.RWMutex
	events  []model.Event
	lastID  uint
	maxSize int
}

func New(maxSize int) (*MemStorage, error) {
	return &MemStorage{
		mu:      sync.RWMutex{},
		maxSize: maxSize,
	}, nil
}

func (s *MemStorage) Close(ctx context.Context) {
	s.events = make([]model.Event, 0)
	s.lastID = 0
	s.maxSize = 0
}

func (s *MemStorage) Create(ctx context.Context, event model.Event) (uint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.lastID++
	event.ID = s.lastID
	s.events = append(s.events, event)

	return event.ID, nil
}

func (s *MemStorage) Update(ctx context.Context, event model.Event) error {
	check := true

	for i, e := range s.events {
		if e.ID == event.ID {
			s.mu.RLock()
			s.events[i] = event
			s.mu.RUnlock()
			check = false
			break
		}
	}

	if check {
		return errors.New("no such model")
	}

	return nil
}

func (s *MemStorage) Delete(ctx context.Context, eventID uint) error {
	findIndex, err := s.findIndexToID(eventID)
	if err != nil {
		return errors.Wrap(err, "cannot delete")
	}

	countEvents := len(s.events)
	if countEvents == 0 {
		return errors.New("cannot delete")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	switch {
	case findIndex == 0:
		if countEvents == 1 {
			s.events = make([]model.Event, 0)
		} else {
			s.events = s.events[1:]
		}
	case findIndex == countEvents-1:
		s.events = s.events[:findIndex-1]
	default:
		s.events = append(s.events[:findIndex-1], s.events[findIndex+1:]...)
	}

	return nil
}

func (s *MemStorage) GetEventByID(ctx context.Context, eventID uint) (model.Event, error) {
	findIndex, err := s.findIndexToID(eventID)
	if err != nil {
		return model.Event{}, errors.Wrap(err, "cannot delete")
	}

	return s.events[findIndex], nil
}

func (s *MemStorage) GetEventsByParams(ctx context.Context, args map[string]interface{}) ([]model.Event, error) {
	var events []model.Event
	countArgs := len(args)

	for _, e := range s.events {
		b := 0
		for s, m := range args {
			if equalFieldEvent(e, s, m) {
				b++
			}
		}

		if b == countArgs {
			events = append(events, e)
		}
	}

	return events, nil
}

func (s *MemStorage) findIndexToID(id uint) (int, error) {
	f := -1

	for i, e := range s.events {
		if e.ID == id {
			f = i
			break
		}
	}

	if f == -1 {
		return 0, errors.New("no such event")
	}

	return f, nil
}

func equalFieldEvent(e model.Event, field string, value interface{}) bool {
	result := false

	switch field {
	case "id":
		if e.ID == value.(uint) {
			result = true
		}
	case "title":
		if e.Title == value.(string) {
			result = true
		}
	case "start_time":
		if e.StartTime == value.(uint) {
			result = true
		}
	case "end_time":
		if e.EndTime == value.(uint) {
			result = true
		}
	case "user_id":
		if e.UserID == value.(uint) {
			result = true
		}
	case "notification_time":
		if e.NotificationTime == value.(uint) {
			result = true
		}
	default:
		result = false
	}

	return result
}
