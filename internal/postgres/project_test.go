package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres"
)

func TestCreateProject(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewProjectService(db)

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME7",
			FirstName: "FIRSTNAME7",
			LastName:  "LASTNAME7",
			Email:     "EMAIL7@email.com",
		})

		project, err := s.CreateProject(
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
