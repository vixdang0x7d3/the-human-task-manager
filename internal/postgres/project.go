package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
)

var _ domain.ProjectService = (*ProjectService)(nil)

type ProjectService struct {
	db     *DB
	logger *logrus.Logger
}

func NewProjectService(db *DB, logger *logrus.Logger) *ProjectService {
	return &ProjectService{
		db:     db,
		logger: logger,
	}
}

func (s *ProjectService) Create(ctx context.Context, title string) (domain.Project, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.Project{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	project, err := createProject(ctx, q, title)
	if err != nil {
		return domain.Project{}, err
	}

	// create owner membership for this user in the project
	if _, err = createOwnerMembership(ctx, q, domain.ProjectMembershipCmd{
		ProjectID: project.ID.String(),
		UserID:    project.UserID.String(),
	}); err != nil {
		return domain.Project(project), err
	}

	return domain.Project(project), nil
}

// ByID return a project with given ID,
// this method is exposed to support
// project find-and-request feature.
// Do not authenticate this method
func (s *ProjectService) ByID(ctx context.Context, id string) (domain.Project, error) {
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

// Find returns projects that this user is the owner,
// sorry for the misguiding name
func (s *ProjectService) Find(ctx context.Context, filter domain.ProjectFilter) ([]domain.Project, int, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return []domain.Project{}, 0, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	projects, n, err := findProjects(ctx, q, filter)
	if err != nil {
		return []domain.Project{}, 0, err
	}

	return Map(projects, toDomainProject), int(n), err
}

func (s *ProjectService) Delete(ctx context.Context, id string) (domain.Project, error) {
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

func findProjects(ctx context.Context, q ProjectQueries, filter domain.ProjectFilter) ([]sqlc.Project, int64, error) {

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return []sqlc.Project{}, 0, &domain.Error{Code: domain.EUNAUTHORIZED, Message: "No user ID in context"}
	}

	rows, err := q.ProjectsByUserID(ctx, sqlc.ProjectsByUserIDParams{
		UserID:  *userID,
		Nlimit:  int32(filter.Limit),
		Noffset: int32(filter.Offset),
	})
	if err != nil {
		return []sqlc.Project{}, 0, err
	}

	if len(rows) == 0 {
		return []sqlc.Project{}, 0, &domain.Error{
			Code:    domain.ENOTFOUND,
			Message: "no projects found",
		}
	}

	return Map(rows, fromProjectsByUserIDRow), rows[0].Count, nil
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

func toDomainProject(project sqlc.Project) domain.Project {
	return domain.Project(project)
}

func fromProjectsByUserIDRow(row sqlc.ProjectsByUserIDRow) sqlc.Project {
	return sqlc.Project{
		ID:     row.ID,
		UserID: row.UserID,
		Title:  row.Title,
	}
}

type ProjectQueries interface {
	ProjectsByUserID(ctx context.Context, arg sqlc.ProjectsByUserIDParams) ([]sqlc.ProjectsByUserIDRow, error)
	ProjectByID(ctx context.Context, id uuid.UUID) (sqlc.Project, error)
	CreateProject(ctx context.Context, arg sqlc.CreateProjectParams) (sqlc.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) (sqlc.Project, error)
}
