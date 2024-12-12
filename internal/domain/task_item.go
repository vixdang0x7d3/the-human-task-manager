package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskItem struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Username        string
	ProjectID       uuid.UUID
	ProjectTitle    string
	CompletedBy     uuid.UUID
	CompletedByName string
	Description     string
	Priority        string
	State           string
	Deadline        time.Time
	Schedule        time.Time
	Wait            time.Time
	Create          time.Time
	End             time.Time
	Tags            []string
	Urgency         float64
}

type TaskItemService interface {
	ByID(ctx context.Context, id string) (TaskItem, error)
	Find(ctx context.Context, filter TaskItemFilter) ([]TaskItem, int, error)
}

type TaskItemFilter struct {
	ProjectID *string

	Q        *string
	State    *string
	Priority *string
	Days     *int64
	Months   *int64

	Offset int
	Limit  int
}
