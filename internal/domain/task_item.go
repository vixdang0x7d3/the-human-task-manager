package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskItem struct {
	ID           uuid.UUID
	Username     string
	ProjectTitle string
	CompletedBy  string
	Description  string
	Priority     string
	State        string
	Deadline     time.Time
	Schedule     time.Time
	Wait         time.Time
	Create       time.Time
	End          time.Time
	Tags         []string
	Urgency      float64
}

type TaskItemService interface {
	TaskItemByID(ctx, id string) (TaskItem, error)
}

type TaskItemFilter struct {
	Search     *string
	Status     *string
	TimePeriod *string

	Offset int
	Limit  int
}
