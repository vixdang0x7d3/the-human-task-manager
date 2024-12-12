package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres"
)

func TestCreateProject(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewProjectService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME7",
			FirstName: "FIRSTNAME7",
			LastName:  "LASTNAME7",
			Email:     "EMAIL7@email.com",
		})

		project, err := s.Create(
			domain.NewContextWithUser(context.Background(), &user),
			"Funny project",
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if reflect.DeepEqual(project.ID, uuid.Nil) {
			t.Errorf("expected ID")
		}

		if !reflect.DeepEqual(project.UserID, user.ID) {
			t.Errorf("user ID mismatch: %q != %q", project.UserID, user.ID)
		}

		if project.Title != "Funny project" {
			t.Errorf("project title mismatch: %q != %q", project.Title, "Funny project")
		}
	})
}

func TestProjectByID(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewProjectService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME8",
			FirstName: "FIRSTNAME8",
			LastName:  "LASTNAME8",
			Email:     "EMAIL8@email.com",
		})

		project := MustCreateProject(
			t,
			domain.NewContextWithUser(context.Background(), &user),
			db, "test project for search by id",
		)

		got, err := s.ByID(context.Background(), project.ID.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got.Title != project.Title {
			t.Errorf("title mismatched: %q != %q", got.Title, project.Title)
		}
	})
}

func TestFindProjects(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewProjectService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME10",
			FirstName: "FIRSTNAME10",
			LastName:  "LASTNAME10",
			Email:     "EMAIL10@email.com",
		})

		MustCreateProject(
			t,
			domain.NewContextWithUser(context.Background(), &user),
			db,
			"test project for search by id",
		)

		MustCreateProject(
			t,
			domain.NewContextWithUser(context.Background(), &user),
			db,
			"test project for search by id",
		)

		MustCreateProject(
			t,
			domain.NewContextWithUser(context.Background(), &user),
			db,
			"test project for search by id",
		)

		projects, n, err := s.Find(
			domain.NewContextWithUser(context.Background(), &user),
			domain.ProjectFilter{Limit: 10, Offset: 0},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(projects) != 3 {
			t.Errorf("expected 3 projects returned, got %d", len(projects))
		}

		if n != len(projects) {
			t.Errorf("also expected 3 total projects in db, got %d", len(projects))
		}
	})
}

func TestDeleteProject(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewProjectService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME9",
			FirstName: "FIRSTNAME9",
			LastName:  "LASTNAME9",
			Email:     "EMAIL9@email.com",
		})

		project := MustCreateProject(
			t,
			domain.NewContextWithUser(context.Background(), &user),
			db,
			"test project for search by id",
		)

		deleted, err := s.Delete(
			domain.NewContextWithUser(context.Background(), &user),
			project.ID.String(),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(deleted.ID, project.ID) {
			t.Errorf("deleted project ID mismatch: %q != %q", deleted.ID, project.ID)
		}

		_, err = s.ByID(context.Background(), project.ID.String())
		if err == nil {
			t.Errorf("expected not found error")
		} else if domain.ErrorCode(err) != domain.ENOTFOUND ||
			domain.ErrorMessage(err) != `projectByID: project ID not found` {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func MustCreateProject(tb testing.TB, ctx context.Context, db *postgres.DB, title string) domain.Project {
	tb.Helper()
	project, err := postgres.NewProjectService(db, logrus.New()).Create(ctx, title)
	if err != nil {
		tb.Fatal(err)
	}
	return project
}
