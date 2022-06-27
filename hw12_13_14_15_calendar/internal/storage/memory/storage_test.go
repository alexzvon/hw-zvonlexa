package memorystorage

import (
	"context"
	"strconv"
	"testing"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/helper"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require"
)

const MAXSIZE = 20

func TestStorageNew(t *testing.T) {
	conn, err := New(MAXSIZE)

	require.Nil(t, err)
	require.NotNil(t, conn)

	conn.Close(context.Background())
}

func TestStorageCreate(t *testing.T) {
	t.Run("Create event", func(t *testing.T) {
		conn, err := New(MAXSIZE)

		require.Nil(t, err)
		require.NotNil(t, conn)

		event := model.Event{
			Title:            "title",
			StartTime:        0,
			EndTime:          0,
			Description:      "description",
			UserID:           12,
			NotificationTime: 17,
		}

		id, err := conn.Create(context.Background(), event)

		require.Nil(t, err)
		require.Equal(t, uint(1), id)

		id, err = conn.Create(context.Background(), event)

		require.Nil(t, err)
		require.Equal(t, uint(2), id)

		id, err = conn.Create(context.Background(), event)

		require.Nil(t, err)
		require.Equal(t, uint(3), id)

		conn.Close(context.Background())
	})
}

func TestStorageGetEventByID(t *testing.T) {
	conn, err := New(MAXSIZE)

	require.Nil(t, err)

	events := []model.Event{
		{
			Title:            "title_1",
			StartTime:        1,
			EndTime:          1,
			Description:      "description_1",
			UserID:           11,
			NotificationTime: 21,
		},
		{
			Title:            "title_2",
			StartTime:        2,
			EndTime:          2,
			Description:      "description_2",
			UserID:           12,
			NotificationTime: 22,
		},
		{
			Title:            "title_3",
			StartTime:        3,
			EndTime:          3,
			Description:      "description_3",
			UserID:           13,
			NotificationTime: 23,
		},
		{
			Title:            "title_4",
			StartTime:        4,
			EndTime:          4,
			Description:      "description_4",
			UserID:           14,
			NotificationTime: 24,
		},
		{
			Title:            "title_5",
			StartTime:        5,
			EndTime:          5,
			Description:      "description_5",
			UserID:           15,
			NotificationTime: 25,
		},
	}

	for i, event := range events {
		event := event

		t.Run(event.Title, func(t *testing.T) {
			id, err := conn.Create(context.Background(), event)

			require.Nil(t, err)
			require.Less(t, uint(0), id)

			events[i].ID = id
		})
	}

	for _, event := range events {
		event := event

		t.Run(helper.ConCat("Get By ID", strconv.Itoa(int(event.ID))), func(t *testing.T) {
			eGet, err := conn.GetEventByID(context.Background(), event.ID)

			require.Nil(t, err)
			require.Equal(t, event, eGet)
		})
	}

	conn.Close(context.Background())
}

func TestStorageUpdate(t *testing.T) {
	t.Run("Update event", func(t *testing.T) {
		conn, err := New(MAXSIZE)

		require.Nil(t, err)
		require.NotNil(t, conn)

		event := model.Event{
			Title:            "title",
			StartTime:        0,
			EndTime:          0,
			Description:      "description",
			UserID:           12,
			NotificationTime: 17,
		}

		id, err := conn.Create(context.Background(), event)

		require.Nil(t, err)
		require.Equal(t, uint(1), id)

		eventUpdate := model.Event{
			ID:               id,
			Title:            "title_update",
			StartTime:        10,
			EndTime:          10,
			Description:      "description_uptade",
			UserID:           12,
			NotificationTime: 17,
		}

		err = conn.Update(context.Background(), eventUpdate)

		require.Nil(t, err)

		eventGet, err := conn.GetEventByID(context.Background(), id)

		require.Nil(t, err)
		require.Equal(t, eventUpdate.Title, eventGet.Title)
		require.Equal(t, eventUpdate.StartTime, eventGet.StartTime)
		require.Equal(t, eventUpdate.EndTime, eventGet.EndTime)
		require.Equal(t, eventUpdate.Description, eventGet.Description)
		require.Equal(t, eventUpdate.UserID, eventGet.UserID)
		require.Equal(t, eventUpdate.NotificationTime, eventGet.NotificationTime)

		conn.Close(context.Background())
	})
}

func TestStorageDelete(t *testing.T) {
	conn, err := New(MAXSIZE)

	require.Nil(t, err)
	require.NotNil(t, conn)

	event := model.Event{
		Title:            "title",
		StartTime:        10,
		EndTime:          15,
		Description:      "description",
		UserID:           11,
		NotificationTime: 5,
	}

	id, err := conn.Create(context.Background(), event)

	require.Nil(t, err)
	require.Less(t, uint(0), id)

	event.ID = id

	eGet, err := conn.GetEventByID(context.Background(), id)

	require.Nil(t, err)
	require.Equal(t, event, eGet)

	err = conn.Delete(context.Background(), id)

	require.Nil(t, err)

	eGet, err = conn.GetEventByID(context.Background(), id)

	require.Error(t, err)
	require.Empty(t, eGet)

	conn.Close(context.Background())
}
