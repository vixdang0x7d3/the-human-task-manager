package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	TaskPriorityH = "H"
	TaskPriorityM = "M"
	TaskPriorityL = "L"
)

type Task struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ProjectID   uuid.UUID
	Description string
	Priority    string
	State       string
	Deadline    time.Time
	Schedule    time.Time
	Wait        time.Time
	Create      time.Time
	End         time.Time
	Tags        []string
}

type TaskService interface {
	Create(ctx context.Context, cmd CreateTaskCmd) (Task, error)
	Delete(ctx context.Context, id string) (Task, error)
	Update(ctx context.Context, id string, cmd UpdateTaskCmd) (Task, error)
	Complete(ctx context.Context, id string) (Task, error)
}

type CreateTaskCmd struct {
	ProjectID   string
	Description string
	Deadline    string
	Schedule    string
	Wait        string
	Priority    string
	Tags        []string
}

type UpdateTaskCmd struct {
	Description string
	Deadline    string
	Schedule    string
	Wait        string
	Priority    string
	Tags        []string
}
