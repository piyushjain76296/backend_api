package models

import "time"

const (
	StatusOpen       = "open"
	StatusInProgress = "in_progress"
	StatusClosed     = "closed"
)

type Ticket struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id,omitempty"` // omitempty if we want to hide it, but req says clear json
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
