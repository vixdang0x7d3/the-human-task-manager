package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
)

type TaskItemService struct {
	db *DB
}

func NewTaskItemService(db *DB) *TaskItemService {
	return &TaskItemService{
		db: db,
	}
}

func (s *TaskItemService) TaskItemByID(ctx context.Context, id string) (domain.TaskItem, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.TaskItem{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	taskItem, err := taskItemByID(ctx, q, id)
	if err != nil {
		return domain.TaskItem{}, err
	}

	return toDomainTaskItem(taskItem), nil
}

func taskItemByID(ctx context.Context, q TaskItemQueries, id string) (sqlc.TaskItem, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return sqlc.TaskItem{}, &domain.Error{Code: domain.EINVALID, Message: "invalid task ID"}
	}

	taskItem, err := q.TaskItemByID(ctx, uuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.TaskItem{}, &domain.Error{Code: domain.ENOTFOUND, Message: "task ID not found"}
		}
	}

	return taskItem, err
}

func toDomainTaskItem(taskItem sqlc.TaskItem) domain.TaskItem {

	var taskPriority string
	if taskItem.Priority != sqlc.TaskPriorityNone {
		taskPriority = string(taskItem.Priority)
	}

	var urgency float64
	if u, ok := numericToFloat64(taskItem.Urgency); ok {
		urgency = u
	}

	return domain.TaskItem{
		ID:           taskItem.ID,
		Username:     taskItem.Username,
		ProjectTitle: taskItem.ProjectTitle,
		CompletedBy:  taskItem.CompletedBy,
		Description:  taskItem.Description,
		Priority:     taskPriority,
		State:        string(taskItem.State),
		Deadline:     taskItem.Deadline,
		Schedule:     taskItem.Schedule,
		Wait:         taskItem.Wait,
		Create:       taskItem.Create,
		End:          taskItem.End,
		Tags:         taskItem.Tags,
		Urgency:      urgency,
	}
}

// FIX: i don't a better way to do this
func numericToFloat64(num pgtype.Numeric) (float64, bool) {
	if !num.Valid {
		return 0, false
	}

	pgFloat, err := num.Float64Value()
	if err != nil || !pgFloat.Valid {
		return 0, false
	}
	return pgFloat.Float64, true
}

type TaskItemQueries interface {
	TaskItemByID(ctx context.Context, id uuid.UUID) (sqlc.TaskItem, error)
}
