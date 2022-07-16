package memorystorage

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/internal/myutils"
	model "github.com/alexzvon/hw-zvonlexa/hw12_13_14_15_calendar/models"
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
		newDT := time.Now()

		conn, err := New(MAXSIZE)

		require.Nil(t, err)
		require.NotNil(t, conn)

		event := model.Event{
			Title:       "title",
			StartDT:     newDT,
			EndDT:       newDT,
			Description: "description",
			UserID:      12,
			NotifDT:     newDT,
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
	newDT := time.Now()
	conn, err := New(MAXSIZE)

	require.Nil(t, err)

	events := []model.Event{
		{
			Title:       "title_1",
			StartDT:     newDT,
			EndDT:       newDT,
			Description: "description_1",
			UserID:      11,
			NotifDT:     newDT,
		},
		{
			Title:       "title_2",
			StartDT:     newDT,
			EndDT:       newDT,
			Description: "description_2",
			UserID:      12,
			NotifDT:     newDT,
		},
		{
			Title:       "title_3",
			StartDT:     newDT,
			EndDT:       newDT,
			Description: "description_3",
			UserID:      13,
			NotifDT:     newDT,
		},
		{
			Title:       "title_4",
			StartDT:     newDT,
			EndDT:       newDT,
			Description: "description_4",
			UserID:      14,
			NotifDT:     newDT,
		},
		{
			Title:       "title_5",
			StartDT:     newDT,
			EndDT:       newDT,
			Description: "description_5",
			UserID:      15,
			NotifDT:     newDT,
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

		t.Run(myutils.ConCat("Get By ID", strconv.Itoa(int(event.ID))), func(t *testing.T) {
			eGet, err := conn.GetEventByID(context.Background(), event.ID)

			require.Nil(t, err)
			require.Equal(t, event, eGet)
		})
	}

	conn.Close(context.Background())
}

func TestStorageUpdate(t *testing.T) {
	t.Run("Update event", func(t *testing.T) {
		newDT := time.Now()
		conn, err := New(MAXSIZE)

		require.Nil(t, err)
		require.NotNil(t, conn)

		event := model.Event{
			Title:       "title",
			StartDT:     newDT,
			EndDT:       newDT,
			Description: "description",
			UserID:      12,
			NotifDT:     newDT,
		}

		id, err := conn.Create(context.Background(), event)

		require.Nil(t, err)
		require.Equal(t, uint(1), id)

		eventUpdate := model.Event{
			ID:          id,
			Title:       "title_update",
			StartDT:     newDT,
			EndDT:       newDT,
			Description: "description_uptade",
			UserID:      12,
			NotifDT:     newDT,
		}

		err = conn.Update(context.Background(), eventUpdate)

		require.Nil(t, err)

		eventGet, err := conn.GetEventByID(context.Background(), id)

		require.Nil(t, err)
		require.Equal(t, eventUpdate.Title, eventGet.Title)
		require.Equal(t, eventUpdate.StartDT, eventGet.StartDT)
		require.Equal(t, eventUpdate.EndDT, eventGet.EndDT)
		require.Equal(t, eventUpdate.Description, eventGet.Description)
		require.Equal(t, eventUpdate.UserID, eventGet.UserID)
		require.Equal(t, eventUpdate.NotifDT, eventGet.NotifDT)

		conn.Close(context.Background())
	})
}

func TestStorageDelete(t *testing.T) {
	newDT := time.Now()
	conn, err := New(MAXSIZE)

	require.Nil(t, err)
	require.NotNil(t, conn)

	event := model.Event{
		Title:       "title",
		StartDT:     newDT,
		EndDT:       newDT,
		Description: "description",
		UserID:      11,
		NotifDT:     newDT,
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
