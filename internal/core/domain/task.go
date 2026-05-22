package domain

import (
	"fmt"
	"time"

	core_errors "github.com/odelshchwank/BigProjectLesson/internal/core/errors"
)

type Task struct {
	ID      int
	Version int

	Title       string
	Description *string

	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserID int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserID int,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: authorUserID,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserID,
	)
}

func (t *Task) Validate() error {
	titleLength := len([]rune(t.Title))
	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf(
			"invalid `Title` len: %d: %w",
			titleLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Description != nil {
		descriptionLength := len([]rune(*t.Description))
		if descriptionLength < 1 || descriptionLength > 1000 {
			return fmt.Errorf(
				"invalid `Description` len: %d: %w",
				descriptionLength,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"`CompletedAt` can't be `nil` if `Completed`==`true`: %w",
				core_errors.ErrInvalidArgument,
			)
		}

		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"`CompletedAt` can't be before `CreatedAt`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"`CompletedAt` must be `nil` if `Completed`==`false`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
