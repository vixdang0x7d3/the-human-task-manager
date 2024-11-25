package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ProjectID   uuid.UUID
	Description string
	Priority    string
	Status      string
	Deadline    time.Time
	Schedule    time.Time
	Wait        time.Time
	Create      time.Time
	End         time.Time
}

type TaskService interface {
	CreateTask(ctx context.Context, cmd CreateTaskCmd) (Task, error)
}

type CreateTaskCmd struct {
	UserID      string
	ProjectID   string
	Description string
	Deadline    string
	Schedule    string
	Wait        string
	Priority    string
}
