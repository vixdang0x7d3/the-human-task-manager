package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
)

type ProjectService struct {
	db *DB
}

func NewProjectService(db *DB) *ProjectService {
	return &ProjectService{
		db: db,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, title string) (domain.Project, error) {
	conn, err := s.db.pool.Acquire(ctx)
	if err != nil {
		return domain.Project{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	project, err := createProject(ctx, q, title)
	if err != nil {
		return domain.Project{}, err
	}

	return domain.Project(project), err
}

func createProject(ctx context.Context, q ProjectQueries, title string) (sqlc.Project, error) {

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return sqlc.Project{}, &domain.Error{Code: domain.EUNAUTHORIZED, Message: "no user ID in context"}
	}

	project, err := q.CreateProject(ctx, sqlc.CreateProjectParams{
		ID:     uuid.New(),
		UserID: *userID,
		Title:  title,
	})
	if err != nil {
		return sqlc.Project{}, err
	}

	return project, nil
}

func deleteProject(ctx context.Context, q ProjectQueries, id string) (sqlc.Project, error)

type ProjectQueries interface {
	CreateProject(ctx context.Context, arg sqlc.CreateProjectParams) (sqlc.Project, error)
}
