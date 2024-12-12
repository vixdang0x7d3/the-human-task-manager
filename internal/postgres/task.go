package postgres

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
)

var _ domain.TaskService = (*TaskService)(nil)

type TaskService struct {
	db     *DB
	logger *logrus.Logger
}

func NewTaskService(db *DB, logger *logrus.Logger) *TaskService {
	return &TaskService{
		db:     db,
		logger: logger,
	}
}

func (s *TaskService) Create(ctx context.Context, cmd domain.CreateTaskCmd) (domain.Task, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	task, err := createTask(ctx, q, cmd)
	if err != nil {
		return toDomainTask(task), err
	}

	return toDomainTask(task), nil
}

func (s *TaskService) Update(ctx context.Context, id string, cmd domain.UpdateTaskCmd) (domain.Task, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return domain.Task{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
	}

	task, err := taskByID(ctx, q, id)
	if err != nil {
		return domain.Task{}, err
	}

	if task.ProjectID.Valid {
		if membership, err := membershipByIDs(ctx, q, domain.ProjectMembershipCmd{
			UserID:    userID.String(),
			ProjectID: task.ProjectID.UUID.String(),
		}); err != nil {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized update, no association with this task",
			}
		} else if membership.Role == sqlc.MembershipRoleInvited || membership.Role == sqlc.MembershipRoleRequested {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized update, unaccepted association with this task",
			}
		}
	} else {
		if *userID != task.UserID {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized update, not task owner",
			}
		}
	}

	task, err = updateTask(ctx, q, task, cmd)
	if err != nil {
		return toDomainTask(task), err
	}

	return toDomainTask(task), nil
}

func (s *TaskService) Complete(ctx context.Context, id string) (domain.Task, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return domain.Task{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
	}

	task, err := taskByID(ctx, q, id)
	if err != nil {
		return domain.Task{}, err
	}

	if task.ProjectID.Valid {
		if membership, err := membershipByIDs(ctx, q, domain.ProjectMembershipCmd{
			UserID:    userID.String(),
			ProjectID: task.ProjectID.UUID.String(),
		}); err != nil {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized update, no association with this task",
			}
		} else if membership.Role == sqlc.MembershipRoleInvited || membership.Role == sqlc.MembershipRoleRequested {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized update, unaccepted association with this task",
			}
		}
	} else {
		if *userID != task.UserID {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized update, not task owner",
			}
		}
	}

	task.CompletedBy = uuid.NullUUID{
		UUID:  *userID,
		Valid: true,
	}
	task, err = completeTask(ctx, q, task)
	if err != nil {
		return toDomainTask(task), err
	}

	return toDomainTask(task), nil
}

func (s *TaskService) Delete(ctx context.Context, id string) (domain.Task, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return domain.Task{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
	}

	task, err := taskByID(ctx, q, id)
	if err != nil {
		return domain.Task{}, err
	}

	if task.ProjectID.Valid {
		if membership, err := membershipByIDs(ctx, q, domain.ProjectMembershipCmd{
			UserID:    userID.String(),
			ProjectID: task.ProjectID.UUID.String(),
		}); err != nil {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized delete, no association with this task",
			}
		} else if membership.Role == sqlc.MembershipRoleInvited || membership.Role == sqlc.MembershipRoleRequested {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized delete, unaccepted association with this task",
			}
		}
	} else {
		if *userID != task.UserID {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized delete, not task owner",
			}
		}
	}

	task, err = deleteTask(ctx, q, task)
	if err != nil {
		return domain.Task{}, err
	}
	return toDomainTask(task), nil
}

// FIX: this is stupid, but i'm too stupid to make this not stupid
func (s *TaskService) SetProject(ctx context.Context, id string, projectID *string) (domain.Task, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return domain.Task{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
	}

	task, err := taskByID(ctx, q, id)
	if err != nil {
		return domain.Task{}, err
	}

	if task.ProjectID.Valid {
		if membership, err := membershipByIDs(ctx, q, domain.ProjectMembershipCmd{
			UserID:    userID.String(),
			ProjectID: task.ProjectID.UUID.String(),
		}); err != nil {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized delete, no association with this task",
			}
		} else if membership.Role == sqlc.MembershipRoleInvited || membership.Role == sqlc.MembershipRoleRequested {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized delete, unaccepted association with this task",
			}
		}
	} else {
		if *userID != task.UserID {
			return domain.Task{}, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "unauthorized delete, not task owner",
			}
		}
	}

	if projectID == nil {
		task, err = setTaskProject(ctx, q, task.ID, uuid.NullUUID{})
		if err != nil {
			return toDomainTask(task), err
		}

		return toDomainTask(task), nil
	}

	// validate the project ID
	membership, err := membershipByIDs(ctx, q, domain.ProjectMembershipCmd{
		UserID:    userID.String(),
		ProjectID: *projectID,
	})
	if err != nil {
		return domain.Task{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "unauthorized update, not a project member",
		}
	} else if membership.Role == sqlc.MembershipRoleInvited || membership.Role == sqlc.MembershipRoleRequested {
		return domain.Task{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "unauthorized update, not a project member yet",
		}

	}

	task, err = setTaskProject(ctx, q, task.ID, uuid.NullUUID{
		UUID:  membership.ProjectID,
		Valid: true,
	})

	return toDomainTask(task), nil
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

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return sqlc.Task{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
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
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "corrupted",
			}

		}
	}

	if cmd.Deadline != "" {
		cmd.Deadline = strings.Join([]string{cmd.Deadline, ":00Z"}, "")
		deadline, err = time.Parse(time.RFC3339, cmd.Deadline)
		if err != nil {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "corrupted deadline timestamp",
			}
		}

		if time.Now().After(deadline) {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "deadline timestamp must be set after current moment",
			}
		}
	}

	if cmd.Schedule != "" {
		cmd.Schedule = strings.Join([]string{cmd.Schedule, ":00Z"}, "")
		schedule, err = time.Parse(time.RFC3339, cmd.Schedule)
		if err != nil {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "corrupted schedule timestamp",
			}
		}

		if time.Now().After(schedule) {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "schedule timestamp must be set after current moment",
			}
		}
	}

	if cmd.Wait != "" {
		cmd.Wait = strings.Join([]string{cmd.Wait, ":00Z"}, "")
		wait, err = time.Parse(time.RFC3339, cmd.Wait)
		if err != nil {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "corrupted wait timestamp"}
		}

		if time.Now().After(wait) {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "wait timestamp must be set after current moment",
			}
		}

		state = sqlc.TaskStateWaiting
	}

	slices.Sort(cmd.Tags)
	task, err := q.CreateTask(ctx, sqlc.CreateTaskParams{
		ID:          uuid.New(),
		UserID:      *userID,
		ProjectID:   projectID,
		Description: cmd.Description,
		Deadline:    deadline,
		Schedule:    schedule,
		Wait:        wait,
		State:       state,
		Priority:    priority,
		Create:      time.Now(),
		Tags:        slices.Compact(cmd.Tags),
	})
	if err != nil {
		return sqlc.Task{}, err
	}

	return task, nil
}

// taskByID find a task with id that matches a given uuid string.
// used in delete and update
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

func updateTask(ctx context.Context, q TaskQueries, task sqlc.Task, cmd domain.UpdateTaskCmd) (sqlc.Task, error) {

	if cmd.Description != "" {
		task.Description = cmd.Description
	}

	if cmd.Deadline != "" {
		cmd.Deadline = strings.Join([]string{cmd.Deadline, ":00Z"}, "")
		deadline, err := time.Parse(time.RFC3339, cmd.Deadline)
		if err != nil {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "corrupted deadline timestamp",
			}
		}

		if time.Now().After(deadline) {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "deadline timestamp must be set after current moment",
			}
		}

		task.Deadline = deadline
	}

	if cmd.Schedule != "" {
		cmd.Schedule = strings.Join([]string{cmd.Schedule, ":00Z"}, "")
		schedule, err := time.Parse(time.RFC3339, cmd.Schedule)
		if err != nil {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "corrupted schedule timestamp",
			}
		}

		if time.Now().After(schedule) {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "schedule timestamp must be set after current moment",
			}
		}

		task.Schedule = schedule
	}

	if cmd.Wait != "" {
		cmd.Wait = strings.Join([]string{cmd.Wait, ":00Z"}, "")
		wait, err := time.Parse(time.RFC3339, cmd.Wait)
		if err != nil {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "corrupted wait timestamp",
			}
		}

		if time.Now().After(wait) {
			return sqlc.Task{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "wait timestamp must be set after current moment",
			}
		}

		task.Wait = wait
	}

	if cmd.Priority != "" {
		var priority sqlc.TaskPriority
		err := priority.Scan(cmd.Priority)
		if err != nil {
			return sqlc.Task{}, err
		}

		task.Priority = priority
	}

	if cmd.Tags != nil {
		slices.Sort(cmd.Tags)
		task.Tags = slices.Compact(cmd.Tags)
	}

	updatedTask, err := q.UpdateTask(ctx, sqlc.UpdateTaskParams{
		ID:          task.ID,
		Description: task.Description,
		Deadline:    task.Deadline,
		Schedule:    task.Schedule,
		Wait:        task.Wait,
		Priority:    task.Priority,
		Tags:        task.Tags,
	})
	if err != nil {
		return updatedTask, err
	}

	return updatedTask, nil
}

func completeTask(ctx context.Context, q TaskQueries, task sqlc.Task) (sqlc.Task, error) {

	completedTask, err := q.CompleteTask(ctx, sqlc.CompleteTaskParams{
		UserID:       task.CompletedBy,
		ID:           task.ID,
		EndTimestamp: time.Now(),
	})
	if err != nil {
		return completedTask, err
	}

	return completedTask, nil
}

func deleteTask(ctx context.Context, q TaskQueries, task sqlc.Task) (sqlc.Task, error) {

	deletedTask, err := q.DeleteTask(ctx, task.ID)

	if err != nil {
		return deletedTask, err
	}

	return deletedTask, nil
}

// FIX: setTaskProject have to accept uuid types instead of string/domain update types
// this breaks the consistent parameter type constraint that other functions comply
// unfortunately, i have no idea on how to fix this garbage code :((
func setTaskProject(ctx context.Context, q TaskQueries, id uuid.UUID, projectID uuid.NullUUID) (sqlc.Task, error) {

	task, err := q.SetTaskProject(ctx, sqlc.SetTaskProjectParams{
		ID:        id,
		ProjectID: projectID,
	})
	if err != nil {
		return task, err
	}

	return task, nil
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
	SetTaskProject(ctx context.Context, arg sqlc.SetTaskProjectParams) (sqlc.Task, error)
}
