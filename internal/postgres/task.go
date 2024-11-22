package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
)

type TaskService struct {
	db *DB
}

func NewTaskService(db *DB) *TaskService {
	return &TaskService{db: db}
}

func (s *TaskService) CreateTask(ctx context.Context, cmd domain.CreateTaskCmd) (domain.Task, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Task{}, errors.New("cannot get a connection")
	}
	q := sqlc.New(conn)

	task, err := createTask(ctx, q, cmd)
	if err != nil {
		return domain.Task{}, err
	}

	return toDomainTask(task), nil
}

func createTask(ctx context.Context, q TaskQueries, cmd domain.CreateTaskCmd) (sqlc.Task, error) {
	panic("unimplemented")
}

func toDomainTask(task sqlc.Task) domain.Task {

	var priority string
	if task.Priority != "none" {
		priority = string(task.Priority)
	}

	var projectID uuid.UUID
	if task.ProjectID.Valid {
		projectID = task.ProjectID.UUID
	}

	return domain.Task{
		ID:          task.ID,
		UserID:      task.UserID,
		ProjectID:   projectID,
		Description: task.Description,
		Status:      string(task.Status),
		Priority:    priority,
		Deadline:    task.Deadline,
		Schedule:    task.Schedule,
		Wait:        task.Wait,
		Create:      task.Create,
		End:         task.End,
	}
}

type TaskQueries interface {
	CreateTask(ctx context.Context, arg sqlc.CreateTaskParams) (sqlc.Task, error)
}
