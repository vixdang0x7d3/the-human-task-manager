package domain

import (
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
	TaskItemByID(ctx, id string) (TaskItem, error)
}

type TaskItemFilter struct {
	Q        *string
	State    *string
	Priority *string
	Days     *int64
	Months   *int64

	Offset int
	Limit  int
}
