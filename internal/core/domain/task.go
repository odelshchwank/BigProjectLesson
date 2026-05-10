package domain

import "time"

type Task struct {
	ID      int
	Version int

	Title       string
	Description *string

	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}
