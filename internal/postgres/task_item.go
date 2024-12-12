package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
)

var _ domain.TaskItemService = (*TaskItemService)(nil)

type TaskItemService struct {
	db     *DB
	logger *logrus.Logger
}

func NewTaskItemService(db *DB, logger *logrus.Logger) *TaskItemService {
	return &TaskItemService{
		db:     db,
		logger: logger,
	}
}

func (s *TaskItemService) ByID(ctx context.Context, id string) (domain.TaskItem, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.TaskItem{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	taskItem, err := taskItemByID(ctx, q, id)
	if err != nil {
		return toDomainTaskItem(taskItem), err
	}

	return toDomainTaskItem(taskItem), nil
}

// Find returns a list of personal tasks that match the given filter
func (s *TaskItemService) Find(ctx context.Context, filter domain.TaskItemFilter) ([]domain.TaskItem, int, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return []domain.TaskItem{}, 0, nil
	}
	defer conn.Release()

	q := sqlc.New(conn)

	if filter.ProjectID != nil {

		userID := domain.UserIDFromContext(ctx)
		if userID == nil {
			return []domain.TaskItem{}, 0, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "Find: no user ID in context",
			}
		}

		if membership, err := membershipByIDs(ctx, q, domain.ProjectMembershipCmd{
			UserID:    userID.String(),
			ProjectID: *filter.ProjectID,
		}); err != nil {
			return []domain.TaskItem{}, 0, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "Find: unauthorized access, not a project member",
			}
		} else if membership.Role == sqlc.MembershipRoleInvited || membership.Role == sqlc.MembershipRoleRequested {
			return []domain.TaskItem{}, 0, &domain.Error{
				Code:    domain.EUNAUTHORIZED,
				Message: "Find: unauthorized access, not a project member yet",
			}
		}

		taskItems, n, err := findTaskItemsByProjectID(ctx, q, filter)
		if err != nil {
			return []domain.TaskItem{}, 0, err
		}
		return Map(taskItems, toDomainTaskItem), int(n), err
	}

	taskItems, n, err := findTaskItemsByUserID(ctx, q, filter)
	if err != nil {
		return []domain.TaskItem{}, 0, err
	}

	return Map(taskItems, toDomainTaskItem), int(n), err
}

func taskItemByID(ctx context.Context, q TaskItemQueries, id string) (sqlc.TaskItem, error) {

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return sqlc.TaskItem{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "taskItemByID: no user ID in context",
		}
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		return sqlc.TaskItem{}, &domain.Error{Code: domain.EINVALID, Message: "invalid task ID"}
	}

	taskItem, err := q.TaskItemByID(ctx, sqlc.TaskItemByIDParams{
		UserID: *userID,
		ID:     uuid,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.TaskItem{}, &domain.Error{Code: domain.ENOTFOUND, Message: "task ID not found"}
		}
	}

	return taskItem, err
}

func findTaskItemsByProjectID(ctx context.Context, q TaskItemQueries, filter domain.TaskItemFilter) ([]sqlc.TaskItem, int64, error) {
	var (
		query        pgtype.Text
		state        sqlc.NullTaskState
		priority     sqlc.NullTaskPriority
		timeInterval pgtype.Interval
	)

	projectID, err := uuid.Parse(*filter.ProjectID)
	if err != nil {
		return []sqlc.TaskItem{}, 0, err
	}

	if filter.Q != nil {
		query = pgtype.Text{
			String: *filter.Q,
			Valid:  true,
		}
	}

	if filter.State != nil {
		err := state.Scan(*filter.State)
		if err != nil {
			return []sqlc.TaskItem{}, 0, &domain.Error{
				Code:    domain.EINVALID,
				Message: "findTaskItemByProjectID: invalid filter for task state",
			}
		}
	}

	if filter.Priority != nil {
		err := priority.Scan(*filter.Priority)
		if err != nil {
			return []sqlc.TaskItem{}, 0, &domain.Error{
				Code:    domain.EINVALID,
				Message: "fidnTaskItemByProjectID: invalid filter for task priority"}
		}
	}

	if filter.Days != nil && filter.Months != nil && *filter.Days > 0 && *filter.Months > 0 {
		var value time.Duration
		if filter.Days != nil {
			value += time.Duration(*filter.Days) * 24 * time.Hour
		}

		if filter.Months != nil {
			value += time.Duration(*filter.Months) * 24 * 30 * time.Hour
		}

		timeInterval = pgtype.Interval{
			Microseconds: value.Microseconds(),
			Valid:        true,
		}
	}

	rows, err := q.FindTaskItemsByProjectID(ctx, sqlc.FindTaskItemsByProjectIDParams{
		ProjectID: uuid.NullUUID{
			UUID:  projectID,
			Valid: true,
		},
		Q:            query,
		State:        state,
		Priority:     priority,
		TimeInterval: timeInterval,
		Nlimit:       int32(filter.Limit),
		Noffset:      int32(filter.Offset),
	})

	if len(rows) == 0 {
		return []sqlc.TaskItem{}, 0, &domain.Error{
			Code:    domain.ENOTFOUND,
			Message: "findTaskItemByProjectID: no tasks found",
		}
	}

	return Map(rows, fromFindTaskItemsByProjectIDRow), rows[0].Count, nil
}

func findTaskItemsByUserID(ctx context.Context, q TaskItemQueries, filter domain.TaskItemFilter) ([]sqlc.TaskItem, int64, error) {
	var (
		query        pgtype.Text
		state        sqlc.NullTaskState
		priority     sqlc.NullTaskPriority
		timeInterval pgtype.Interval
	)

	if filter.Q != nil {
		query = pgtype.Text{
			String: *filter.Q,
			Valid:  true,
		}
	}

	if filter.State != nil {
		err := state.Scan(*filter.State)
		if err != nil {
			return []sqlc.TaskItem{}, 0, &domain.Error{
				Code:    domain.EINVALID,
				Message: "findTaskItemByUserID: invalid filter for task state"}
		}
	}

	if filter.Priority != nil {
		err := priority.Scan(*filter.Priority)
		if err != nil {
			return []sqlc.TaskItem{}, 0, &domain.Error{
				Code:    domain.EINVALID,
				Message: "findTaskItemByUserID: invalid filter for task priorirty",
			}
		}
	}

	if filter.Days != nil && filter.Months != nil && *filter.Days > 0 && *filter.Months > 0 {
		var value time.Duration
		if filter.Days != nil {
			value += time.Duration(*filter.Days) * 24 * time.Hour
		}

		if filter.Months != nil {
			value += time.Duration(*filter.Months) * 24 * 30 * time.Hour
		}

		timeInterval = pgtype.Interval{
			Microseconds: value.Microseconds(),
			Valid:        true,
		}
	}

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return []sqlc.TaskItem{}, 0, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "findTaskItemByUserID: no user ID from context",
		}
	}

	rows, err := q.FindTaskItemsByUserID(ctx, sqlc.FindTaskItemsByUserIDParams{
		UserID:       *userID,
		Q:            query,
		State:        state,
		Priority:     priority,
		TimeInterval: timeInterval,
		Nlimit:       int32(filter.Limit),
		Noffset:      int32(filter.Offset),
	})
	if err != nil {
		return []sqlc.TaskItem{}, 0, err
	}

	// prevents indexing an empty array
	if len(rows) == 0 {
		return []sqlc.TaskItem{}, 0, &domain.Error{
			Code:    domain.ENOTFOUND,
			Message: "findTaskItemsByUserID: no tasks found",
		}
	}

	return Map(rows, fromFindTaskItemsByUserIDRow), rows[0].Count, err
}

func toDomainTaskItem(taskItem sqlc.TaskItem) domain.TaskItem {

	var projectID uuid.UUID
	if taskItem.ProjectID.Valid {
		projectID = taskItem.ProjectID.UUID
	}

	var completedBy uuid.UUID
	if taskItem.CompletedBy.Valid {
		completedBy = taskItem.CompletedBy.UUID
	}

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
		UserID:       taskItem.UserID,
		Username:     taskItem.Username,
		ProjectID:    projectID,
		ProjectTitle: taskItem.ProjectTitle,
		CompletedBy:  completedBy,
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

func fromFindTaskItemsByProjectIDRow(row sqlc.FindTaskItemsByProjectIDRow) sqlc.TaskItem {
	return sqlc.TaskItem{
		ID:              row.ID,
		UserID:          row.UserID,
		Username:        row.Username,
		ProjectID:       row.ProjectID,
		ProjectTitle:    row.ProjectTitle,
		CompletedBy:     row.CompletedBy,
		CompletedByName: row.CompletedByName,
		Description:     row.Description,
		Priority:        row.Priority,
		State:           row.State,
		Deadline:        row.Deadline,
		Schedule:        row.Schedule,
		Wait:            row.Wait,
		Create:          row.Create,
		End:             row.End,
		Tags:            row.Tags,
		Urgency:         row.Urgency,
	}
}

func fromFindTaskItemsByUserIDRow(row sqlc.FindTaskItemsByUserIDRow) sqlc.TaskItem {
	return sqlc.TaskItem{
		ID:              row.ID,
		UserID:          row.UserID,
		Username:        row.Username,
		ProjectID:       row.ProjectID,
		ProjectTitle:    row.ProjectTitle,
		CompletedBy:     row.CompletedBy,
		CompletedByName: row.CompletedByName,
		Description:     row.Description,
		Priority:        row.Priority,
		State:           row.State,
		Deadline:        row.Deadline,
		Schedule:        row.Schedule,
		Wait:            row.Wait,
		Create:          row.Create,
		End:             row.End,
		Tags:            row.Tags,
		Urgency:         row.Urgency,
	}

}

type TaskItemQueries interface {
	TaskItemByID(ctx context.Context, arg sqlc.TaskItemByIDParams) (sqlc.TaskItem, error)
	FindTaskItemsByUserID(ctx context.Context, arg sqlc.FindTaskItemsByUserIDParams) ([]sqlc.FindTaskItemsByUserIDRow, error)
	FindTaskItemsByProjectID(ctx context.Context, arg sqlc.FindTaskItemsByProjectIDParams) ([]sqlc.FindTaskItemsByProjectIDRow, error)
}
