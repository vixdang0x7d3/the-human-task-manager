package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
)

var _ domain.TaskService = (*TaskService)(nil)

type TaskService struct {
	db *DB
}

func NewTaskService(db *DB) *TaskService {
	return &TaskService{db: db}
}

func (s *TaskService) CreateTask(ctx context.Context, cmd domain.CreateTaskCmd) (domain.Task, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	task, err := createTask(ctx, q, cmd)
	if err != nil {
		return domain.Task{}, err
	}

	return toDomainTask(task), nil
}

func createTask(ctx context.Context, q TaskQueries, cmd domain.CreateTaskCmd) (sqlc.Task, error) {

	var (
		err       error
		projectID uuid.NullUUID
		deadline  time.Time
		schedule  time.Time
		wait      time.Time
		status    sqlc.TaskStatus
		priority  sqlc.TaskPriority
	)

	status = sqlc.TaskStatusStarted
	priority = sqlc.TaskPriorityNone

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return sqlc.Task{}, &domain.Error{Code: domain.EINVALID, Message: "corrupted user ID"}

	}

	if cmd.ProjectID != "" {

		val, err := uuid.Parse(cmd.ProjectID)
		if err != nil {
			return sqlc.Task{}, &domain.Error{Code: domain.EINVALID, Message: "corrupted project ID"}
		}

		projectID = uuid.NullUUID{
			UUID:  val,
			Valid: true,
		}
	}

	if cmd.Priority != "" {
		if err = priority.Scan(cmd.Priority); err != nil {
			return sqlc.Task{}, &domain.Error{Code: domain.EINVALID, Message: "corrupted priority"}
		}
	}

	if cmd.Deadline != "" {
		cmd.Deadline = strings.Join([]string{cmd.Deadline, ":00Z"}, "")
		deadline, err = time.Parse(time.RFC3339, cmd.Deadline)
		if err != nil {
			return sqlc.Task{}, &domain.Error{Code: domain.EINVALID, Message: "corrupted deadline timestamp"}
		}
	}

	if cmd.Schedule != "" {
		cmd.Schedule = strings.Join([]string{cmd.Schedule, ":00Z"}, "")
		schedule, err = time.Parse(time.RFC3339, cmd.Schedule)
		if err != nil {
			return sqlc.Task{}, &domain.Error{Code: domain.EINVALID, Message: "corrupted schedule timestamp"}
		}
	}

	if cmd.Wait != "" {
		cmd.Wait = strings.Join([]string{cmd.Wait, ":00Z"}, "")
		wait, err = time.Parse(time.RFC3339, cmd.Wait)
		if err != nil {
			return sqlc.Task{}, &domain.Error{Code: domain.EINVALID, Message: "corrupted wait timestamp"}
		}
		status = sqlc.TaskStatusWaiting
	}

	task, err := q.CreateTask(ctx, sqlc.CreateTaskParams{
		ID:          uuid.New(),
		UserID:      userID,
		ProjectID:   projectID,
		Description: cmd.Description,
		Deadline:    deadline,
		Schedule:    schedule,
		Wait:        wait,
		Status:      status,
		Priority:    priority,
		Create:      time.Now(),
	})
	if err != nil {
		return sqlc.Task{}, err
	}

	return task, nil
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
