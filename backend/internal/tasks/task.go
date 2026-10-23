// Package tasks contains the business logic of student study tasks.
package tasks

import (
	"errors"
	"strings"
	"time"
)

// Status is the progress status of a task.
type Status string

const (
	// StatusPlanned is a task that has not been started.
	StatusPlanned Status = "planned"
	// StatusInProgress is a task that is being worked on.
	StatusInProgress Status = "in_progress"
	// StatusDone is a finished task.
	StatusDone Status = "done"
)

var (
	// ErrEmptyTitle is returned when the title is empty.
	ErrEmptyTitle = errors.New("title must not be empty")
	// ErrDeadlineInPast is returned when the deadline is already in the past.
	ErrDeadlineInPast = errors.New("deadline must not be in the past")
	// ErrInvalidStatus is returned for an unknown status.
	ErrInvalidStatus = errors.New("invalid task status")
)

// Task is a study task with a deadline.
type Task struct {
	ID          int64
	Title       string
	Description string
	Deadline    time.Time
	Status      Status
	CreatedAt   time.Time
}

// Valid reports whether s is one of the known statuses.
func (s Status) Valid() bool {
	switch s {
	case StatusPlanned, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

// New validates the input and creates a task in the planned status.
func New(title, description string, deadline, now time.Time) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, ErrEmptyTitle
	}
	if deadline.Before(now) {
		return Task{}, ErrDeadlineInPast
	}
	return Task{
		Title:       title,
		Description: description,
		Deadline:    deadline,
		Status:      StatusPlanned,
		CreatedAt:   now,
	}, nil
}

// ChangeStatus switches the task to a new valid status.
func (t *Task) ChangeStatus(next Status) error {
	if !next.Valid() {
		return ErrInvalidStatus
	}
	t.Status = next
	return nil
}
