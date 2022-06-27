package model

type Event struct {
	ID               uint   `json:"id" db:"id"`
	Title            string `json:"title" db:"title"`
	StartTime        uint   `json:"startTime" db:"start_time"`
	EndTime          uint   `json:"endTime" db:"end_time"`
	Description      string `json:"description" db:"description"`
	UserID           uint   `json:"userId" db:"user_id"`
	NotificationTime uint   `json:"notificationTime" db:"notification_time"`
}
