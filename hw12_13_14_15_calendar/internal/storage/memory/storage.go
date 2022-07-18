package memorystorage

import (
	"context"
	"sync"
	"time"

	model "github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/models"
	"github.com/pkg/errors"
)

type MemStorage struct {
	mu      sync.RWMutex
	events  map[uint]model.Event
	lastID  uint
	maxSize int
}

func New(maxSize int) (*MemStorage, error) {
	return &MemStorage{
		mu:      sync.RWMutex{},
		events:  make(map[uint]model.Event),
		maxSize: maxSize,
	}, nil
}

func (s *MemStorage) Close(ctx context.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.events = make(map[uint]model.Event, 0)
	s.lastID = 0
	s.maxSize = 0
}

func (s *MemStorage) Create(ctx context.Context, event model.Event) (uint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.lastID++
	event.ID = s.lastID
	s.events[event.ID] = event

	return event.ID, nil
}

func (s *MemStorage) Update(ctx context.Context, event model.Event) error {
	if _, ok := s.events[event.ID]; ok {
		s.mu.RLock()
		s.events[event.ID] = event
		s.mu.RUnlock()

		return nil
	}

	return errors.New("no such model")
}

func (s *MemStorage) Delete(ctx context.Context, eventID uint) error {
	if _, ok := s.events[eventID]; ok {
		s.mu.RLock()
		delete(s.events, eventID)
		s.mu.RUnlock()

		return nil
	}

	return errors.New("no such model")
}

func (s *MemStorage) GetEventByID(ctx context.Context, eventID uint) (model.Event, error) {
	if _, ok := s.events[eventID]; ok {
		return s.events[eventID], nil
	}

	return model.Event{}, errors.New("no such model")
}

func (s *MemStorage) GetEventsByParams(ctx context.Context, fields map[string]interface{}) ([]model.Event, error) {
	var events []model.Event
	countFields := len(fields)

	for _, event := range s.events {
		countEquelField := 0
		for nameField, valueField := range fields {
			if equalFieldEvent(event, nameField, valueField) {
				countEquelField++
			}
		}

		if countEquelField == countFields {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *MemStorage) ListEventsToDay(ctx context.Context, dt time.Time) ([]model.Event, error) {
	var events []model.Event

	for _, event := range s.events {
		if equalFieldEvent(event, "start_dt", dt) {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *MemStorage) ListEventsToWeek(ctx context.Context, dt time.Time) ([]model.Event, error) {
	var events []model.Event
	et := dt.Add(7)

	for _, event := range s.events {
		if dt.Before(event.StartDT) && et.After(event.StartDT) {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *MemStorage) ListEventsToMonth(ctx context.Context, dt time.Time) ([]model.Event, error) {
	var events []model.Event
	et := dt.Add(30)

	for _, event := range s.events {
		if dt.Before(event.StartDT) && et.After(event.StartDT) {
			events = append(events, event)
		}
	}

	return events, nil
}

func equalFieldEvent(event model.Event, field string, value interface{}) bool {
	result := false

	switch field {
	case "id":
		if event.ID == value.(uint) {
			result = true
		}
	case "title":
		if event.Title == value.(string) {
			result = true
		}
	case "start_dt":
		if event.StartDT == value.(time.Time) {
			result = true
		}
	case "end_dt":
		if event.EndDT == value.(time.Time) {
			result = true
		}
	case "user_id":
		if event.UserID == value.(uint) {
			result = true
		}
	case "notif_dt":
		if event.NotifDT == value.(time.Time) {
			result = true
		}
	default:
		result = false
	}

	return result
}
