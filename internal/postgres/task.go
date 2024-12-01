package postgres

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (s *TaskService) StartWaitingTaskWorker(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tasks, err := s.StartWaitingTasks(ctx)
			if err != nil {
				log.Printf("Failed to start waiting tasks: %v", err)
				continue
			}

			log.Printf("Started %d waiting tasks\n", len(tasks))
			for _, t := range tasks {
				log.Printf("Waiting task %q started...", t.ID)
			}
		}
	}
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

func (s *TaskService) UpdateTask(ctx context.Context, id string, cmd domain.UpdateTaskCmd) (domain.Task, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	task, err := updateTask(ctx, id, q, cmd)
	if err != nil {
		return domain.Task{}, err
	}

	return toDomainTask(task), nil
}

func (s *TaskService) CompleteTask(ctx context.Context, id string) (domain.Task, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	task, err := completeTask(ctx, q, id)
	if err != nil {
		return domain.Task{}, err
	}

	return toDomainTask(task), nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) (domain.Task, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	task, err := deleteTask(ctx, q, id)
	if err != nil {
		return domain.Task{}, err
	}
	return toDomainTask(task), nil
}

func (s *TaskService) StartWaitingTasks(ctx context.Context) ([]domain.Task, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return []domain.Task{}, err
	}
	defer tx.Rollback(ctx)

	q := sqlc.New(tx)
	started, err := startWaitingTask(ctx, q)
	if err != nil {
		return []domain.Task{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return []domain.Task{}, err
	}

	return Map(started, toDomainTask), nil
}

func createTask(ctx context.Context, q TaskQueries, cmd domain.CreateTaskCmd) (sqlc.Task, error) {

	var (
		err       error
		projectID uuid.NullUUID
		deadline  time.Time
		schedule  time.Time
		wait      time.Time
		state     sqlc.TaskState
		priority  sqlc.TaskPriority
	)

	state = sqlc.TaskStateStarted
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
		state = sqlc.TaskStateWaiting
	}

	task, err := q.CreateTask(ctx, sqlc.CreateTaskParams{
		ID:          uuid.New(),
		UserID:      userID,
		ProjectID:   projectID,
		Description: cmd.Description,
		Deadline:    deadline,
		Schedule:    schedule,
		Wait:        wait,
		State:       state,
		Priority:    priority,
		Create:      time.Now(),
		Tags:        cmd.Tags,
	})
	if err != nil {
		return sqlc.Task{}, err
	}

	return task, nil
}

func taskByID(ctx context.Context, q TaskQueries, id string) (sqlc.Task, error) {

	uuid, err := uuid.Parse(id)
	if err != nil {
		return sqlc.Task{}, &domain.Error{Code: domain.EINVALID, Message: "corrupted task ID"}
	}

	task, err := q.TaskByID(ctx, uuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.Task{}, &domain.Error{Code: domain.ENOTFOUND, Message: "task ID not found"}
		}
		return sqlc.Task{}, err
	}

	return task, nil
}

// TODO: task update authorization:
// compare user ID from context with user ID from TaskByID,
// if not equal, don't allow update.
func updateTask(ctx context.Context, id string, q TaskQueries, cmd domain.UpdateTaskCmd) (sqlc.Task, error) {
	task, err := taskByID(ctx, q, id)
	if err != nil {
		return task, err
	}

	if cmd.Description != "" {
		task.Description = cmd.Description
	}

	if cmd.Deadline != "" {
		cmd.Deadline = strings.Join([]string{cmd.Deadline, ":00Z"}, "")
		deadline, err := time.Parse(time.RFC3339, cmd.Deadline)
		if err != nil {
			// do logging here idk
		} else {
			task.Deadline = deadline
		}
	}

	if cmd.Schedule != "" {
		cmd.Schedule = strings.Join([]string{cmd.Schedule, ":00Z"}, "")
		schedule, err := time.Parse(time.RFC3339, cmd.Schedule)
		if err != nil {
			// do logging here idk
		} else {
			task.Schedule = schedule
		}
	}

	if cmd.Wait != "" {
		cmd.Wait = strings.Join([]string{cmd.Wait, ":00Z"}, "")
		wait, err := time.Parse(time.RFC3339, cmd.Wait)
		if err != nil {
			// do logging here idk
		} else {
			task.Wait = wait
		}
	}

	if cmd.Priority != "" {
		var priority sqlc.TaskPriority
		if err := priority.Scan(cmd.Priority); err != nil {
			// do some logging here
		} else {
			task.Priority = priority
		}
	}

	if len(cmd.Tags) != 0 {
		task.Tags = cmd.Tags
	}

	updatedTask, err := q.UpdateTask(ctx, sqlc.UpdateTaskParams{
		ID:          task.ID,
		Description: task.Description,
		Priority:    task.Priority,
		Deadline:    task.Deadline,
		Schedule:    task.Schedule,
		Wait:        task.Wait,
		Tags:        task.Tags,
	})
	if err != nil {
		return updatedTask, err
	}

	return updatedTask, nil
}

// TODO: task complete authorization:
// compare user ID from context with user ID from TaskByID,
// if not equal, don't allow update.
func completeTask(ctx context.Context, q TaskQueries, id string) (sqlc.Task, error) {
	task, err := taskByID(ctx, q, id)
	if err != nil {
		return task, err
	}

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return task, &domain.Error{Code: domain.EUNAUTHORIZED, Message: "no user ID from context"}
	}

	completedTask, err := q.CompleteTask(ctx, sqlc.CompleteTaskParams{
		UserID: uuid.NullUUID{
			UUID:  *userID,
			Valid: true,
		},
		ID: task.ID,
	})
	if err != nil {
		return completedTask, err
	}

	return completedTask, nil
}

func deleteTask(ctx context.Context, q TaskQueries, id string) (sqlc.Task, error) {
	task, err := taskByID(ctx, q, id)
	if err != nil {
		return task, err
	}

	deletedTask, err := q.DeleteTask(ctx, task.ID)
	if err != nil {
		return deletedTask, err
	}

	return deletedTask, nil
}

func startWaitingTask(ctx context.Context, q TaskQueries) ([]sqlc.Task, error) {
	tasks, err := q.StartWaitingTasks(ctx)
	if err != nil {
		return []sqlc.Task{}, err
	}

	return tasks, nil
}

func toDomainTask(task sqlc.Task) domain.Task {

	var priority string
	if task.Priority != sqlc.TaskPriorityNone {
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
		State:       string(task.State),
		Priority:    priority,
		Deadline:    task.Deadline,
		Schedule:    task.Schedule,
		Wait:        task.Wait,
		Create:      task.Create,
		End:         task.End,
		Tags:        task.Tags,
	}
}

type TaskQueries interface {
	TaskByID(ctx context.Context, id uuid.UUID) (sqlc.Task, error)
	CreateTask(ctx context.Context, arg sqlc.CreateTaskParams) (sqlc.Task, error)
	UpdateTask(ctx context.Context, arg sqlc.UpdateTaskParams) (sqlc.Task, error)
	CompleteTask(ctx context.Context, arg sqlc.CompleteTaskParams) (sqlc.Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) (sqlc.Task, error)
	StartWaitingTasks(ctx context.Context) ([]sqlc.Task, error)
}
