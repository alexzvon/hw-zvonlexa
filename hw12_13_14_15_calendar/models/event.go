package model

import "time"

type Event struct {
	ID          uint      `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	StartDT     time.Time `json:"startTime" db:"start_dt"`
	EndDT       time.Time `json:"endTime" db:"end_dt"`
	Description string    `json:"description" db:"description"`
	UserID      uint      `json:"userId" db:"user_id"`
	NotifDT     time.Time `json:"notificationTime" db:"notif_dt"`
}
