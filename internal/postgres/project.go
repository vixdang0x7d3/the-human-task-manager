package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	return domain.Project(project), nil
}

func (s *ProjectService) ProjectByID(ctx context.Context, id string) (domain.Project, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Project{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	project, err := projectByID(ctx, q, id)
	if err != nil {
		return domain.Project{}, err
	}
	return domain.Project(project), nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, id string) (domain.Project, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Project{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	project, err := deleteProject(ctx, q, id)
	if err != nil {
		return domain.Project{}, err
	}
	return domain.Project(project), nil
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

func projectByID(ctx context.Context, q ProjectQueries, id string) (sqlc.Project, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return sqlc.Project{}, err
	}
	project, err := q.ProjectByID(ctx, uuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.Project{}, &domain.Error{Code: domain.ENOTFOUND, Message: "project ID not found"}
		}

		return sqlc.Project{}, err
	}
	return project, nil
}

func deleteProject(ctx context.Context, q ProjectQueries, id string) (sqlc.Project, error) {
	project, err := projectByID(ctx, q, id)
	if err != nil {
		return sqlc.Project{}, err
	}

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return sqlc.Project{}, &domain.Error{Code: domain.EUNAUTHORIZED, Message: "no user ID in context"}
	}

	if *userID != project.UserID {
		return sqlc.Project{}, &domain.Error{Code: domain.EUNAUTHORIZED, Message: "unauthorized user, cannot delete"}
	}

	deleted, err := q.DeleteProject(ctx, project.ID)
	if err != nil {
		return sqlc.Project{}, err
	}

	return deleted, nil
}

type ProjectQueries interface {
	ProjectByID(ctx context.Context, id uuid.UUID) (sqlc.Project, error)
	CreateProject(ctx context.Context, arg sqlc.CreateProjectParams) (sqlc.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) (sqlc.Project, error)
}
