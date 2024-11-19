package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
)

type CreateTaskCmd struct {
	UserID      string
	ProjectID   string
	Description string
	Deadline    string
	Schedule    string
	Wait        string
	Priority    string
}

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

func ToDomainTask(task database.Task) Task {

	var priority string
	if task.Priority != database.TaskPriorityNone {
		priority = string(task.Priority)
	}

	var projectID uuid.UUID
	if task.ProjectID.Valid {
		projectID = task.ProjectID.UUID
	}

	return Task{
		ID:          task.ID,
		UserID:      task.UserID,
		ProjectID:   projectID,
		Description: task.Description,
		Priority:    priority,
		Status:      string(task.Status),
		Deadline:    task.Deadline,
		Schedule:    task.Schedule,
		Wait:        task.Wait,
		Create:      task.Create,
		End:         task.End,
	}
}
