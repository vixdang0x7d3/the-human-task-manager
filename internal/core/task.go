package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/types"
)

type TaskCore struct {
	Store TaskStore
}

func NewTaskCore(store TaskStore) *TaskCore {
	return &TaskCore{
		Store: store,
	}
}

func (c *TaskCore) CreateTask(ctx context.Context, arg types.CreateTaskCmd) (types.Task, error) {

	var (
		err       error
		projectID uuid.UUID
		deadline  time.Time
		schedule  time.Time
		wait      time.Time
		priority  database.TaskPriority
		status    database.TaskStatus
	)

	status = database.TaskStatusStarted

	userID, err := uuid.Parse(arg.UserID)
	if err != nil {
		return types.Task{}, fmt.Errorf("create task error: %w", err)
	}

	if arg.ProjectID != "" {
		projectID, err = uuid.Parse(arg.ProjectID)
		if err != nil {
			return types.Task{}, fmt.Errorf("create task error: %w", err)
		}
	}

	priority = database.TaskPriorityNone
	if arg.Priority != "" {
		err := priority.Scan(arg.Priority)
		if err != nil {
			return types.Task{}, err
		}
	}

	if arg.Deadline != "" {

		arg.Deadline = strings.Join([]string{arg.Deadline, ":00Z"}, "")
		deadline, err = time.Parse(time.RFC3339, arg.Deadline)
		if err != nil {
			return types.Task{}, err
		}

		if deadline.Compare(time.Now()) <= 0 {
			return types.Task{}, errors.New("deadline with invalid timestamp")
		}
	}

	if arg.Schedule != "" {
		arg.Schedule = strings.Join([]string{arg.Schedule, ":00Z"}, "")
		schedule, err = time.Parse(time.RFC3339, arg.Schedule)
		if err != nil {
			return types.Task{}, err
		}
	}

	if arg.Wait != "" {
		arg.Wait = strings.Join([]string{arg.Wait, ":00Z"}, "")
		wait, err = time.Parse(time.RFC3339, arg.Wait)
		if err != nil {
			return types.Task{}, err
		}

		status = database.TaskStatusWaiting
	}

	task, err := c.Store.CreateTask(ctx, database.CreateTaskParams{
		ID: uuid.New(),
		ProjectID: uuid.NullUUID{
			UUID:  projectID,
			Valid: true,
		},
		UserID:      userID,
		Description: arg.Description,
		Deadline:    deadline,
		Schedule:    schedule,
		Wait:        wait,
		Create:      time.Now(),
		End:         time.Time{},
		Priority:    priority,
		Status:      status,
	})

	if err != nil {
		return types.Task{}, err
	}

	return types.ToDomainTask(task), nil
}

type TaskStore interface {
	CreateTask(ctx context.Context, arg database.CreateTaskParams) (database.Task, error)
}
