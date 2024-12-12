package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
)

func TestCreateTask(t *testing.T) {

	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewTaskService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {

		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "chineseman",
			FirstName: "Franz",
			LastName:  "Kafka",
			Email:     "kafka.franz@email.com",
		})

		ctxWithUser := domain.NewContextWithUser(context.Background(), &user)

		cmd := domain.CreateTaskCmd{
			ProjectID:   "",
			Description: "DESCRIPTION0",
			Deadline:    "2024-12-23T07:00",
			Schedule:    "",
			Wait:        "",
			Priority:    "",
			Tags:        []string{"TAG1", "TAG2", "TAG3"},
		}

		task, err := s.Create(ctxWithUser, cmd)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if reflect.DeepEqual(task.ID, uuid.Nil) {
			t.Errorf("expected id")
		}

		if !reflect.DeepEqual(task.UserID, user.ID) {
			t.Errorf("user ID mismatch: %q != %q", task.UserID, user.ID)
		}

		if task.Description != cmd.Description {
			t.Errorf("description mismatch: %q != %q", task.Description, cmd.Description)
		}

		wantDeadline, _ := time.Parse(time.RFC3339, "2024-12-23T07:00:00Z")
		if task.Deadline.Compare(wantDeadline) != 0 {
			t.Errorf("deadline mismatch: %v = %v", task.Deadline, wantDeadline)
		}

		if task.Create.IsZero() {
			t.Errorf("expected create timestamp")
		}

		if task.State != string(sqlc.TaskStateStarted) {
			t.Errorf("expected status: %q != %q", task.State, string(sqlc.TaskStateStarted))
		}

		if !reflect.DeepEqual(cmd.Tags, task.Tags) {
			t.Errorf("tags mismatched %#v != %#v", cmd.Tags, task.Tags)
		}
	})
}

func TestUpdateTask(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	s := postgres.NewTaskService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME2",
			FirstName: "FIRSTNAME2",
			LastName:  "LASTNAME2",
			Email:     "EMAIL2@email.com",
			Password:  "fckthisshit",
		})

		ctxWithUser := domain.NewContextWithUser(context.Background(), &user)

		task := MustCreateTask(t, ctxWithUser, db, domain.CreateTaskCmd{
			ProjectID:   "",
			Description: "DESCRIPTION1",
			Deadline:    "2024-12-13T07:00",
			Schedule:    "2024-12-13T07:00",
			Wait:        "",
			Priority:    "M",
			Tags:        []string{"school", "INT10111"},
		})

		s.Update(ctxWithUser, task.ID.String(), domain.UpdateTaskCmd{
			Description: "ModifiedDescription",
			Deadline:    "",
			Schedule:    "",
			Wait:        "",
			Priority:    "H",
			Tags:        []string{"school"},
		})

		taskItem := MustTaskItemByID(t, ctxWithUser, db, task.ID.String())

		if taskItem.Description != "ModifiedDescription" {
			t.Errorf("expected new description: %q != %q", taskItem.Description, "ModifiedDescription")
		}

		if taskItem.Priority != "H" {
			t.Errorf("expected new priority: %q != %q", taskItem.Priority, "H")
		}

		if !reflect.DeepEqual(taskItem.Tags, []string{"school"}) {
			t.Errorf("expected new tags list: %v != %v", taskItem.Tags, []string{"school"})
		}
	})
}

func TestCompleteTask(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewTaskService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME4",
			FirstName: "FIRSTNAME4",
			LastName:  "LASTNAME4",
			Email:     "EMAIL4@email.com",
			Password:  "llalalalala",
		})

		ctxWithUser := domain.NewContextWithUser(context.Background(), &user)

		task := MustCreateTask(t, ctxWithUser, db, domain.CreateTaskCmd{
			ProjectID:   "",
			Description: "DESCRIPTION4",
			Deadline:    "2024-12-13T07:00",
			Schedule:    "2024-12-13T07:00",
			Wait:        "",
			Priority:    "L",
			Tags:        []string{"A", "B"},
		})

		_, err := s.Complete(ctxWithUser, task.ID.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		completed := MustTaskItemByID(t, ctxWithUser, db, task.ID.String())
		if completed.State != "completed" {
			t.Errorf("wrong state: %q != %q", completed.State, "completed")
		}

		if !reflect.DeepEqual(completed.CompletedBy, user.ID) {
			t.Errorf("expected complete_by, got %q want %q", completed.CompletedBy, user.ID)
		}

		if completed.End.IsZero() {
			t.Errorf("expected end timestamp")
		}
	})
}

func TestDeleteTask(t *testing.T) {

	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewTaskService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME5",
			FirstName: "FIRSTNAME5",
			LastName:  "LASTNAME5",
			Email:     "EMAIL5@email.com",
			Password:  "bababababa",
		})

		ctxWithUser := domain.NewContextWithUser(context.Background(), &user)

		task := MustCreateTask(t, ctxWithUser, db, domain.CreateTaskCmd{
			ProjectID:   "",
			Description: "DESCRIPTION5",
			Deadline:    "2024-12-13T07:00",
			Schedule:    "2024-12-13T07:00",
			Wait:        "",
			Priority:    "L",
			Tags:        []string{"B", "C"},
		})

		_, err := s.Delete(ctxWithUser, task.ID.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		deleted := MustTaskItemByID(t, ctxWithUser, db, task.ID.String())
		if deleted.State != "deleted" {
			t.Errorf("wrong state: %q != %q", deleted.State, "deleted")
		}
	})
}

func TestSetProject(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewTaskService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {

		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME19",
			FirstName: "FIRSTNAME19",
			LastName:  "LASTNAME19",
			Email:     "EMAIL19@email.com",
			Password:  "bababababa",
		})

		ctxWithUser := domain.NewContextWithUser(context.Background(), &user)

		task := MustCreateTask(t, ctxWithUser, db, domain.CreateTaskCmd{
			ProjectID:   "",
			Description: "DESCRIPTION6",
			Deadline:    "2024-12-13T07:00",
			Schedule:    "2024-12-13T07:00",
			Wait:        "",
			Priority:    "L",
			Tags:        []string{"B", "C"},
		})

		project := MustCreateProject(t, ctxWithUser, db, "PROJECT123")
		projectID := project.ID.String()

		_, err := s.SetProject(ctxWithUser, task.ID.String(), &projectID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		updated := MustTaskItemByID(t, ctxWithUser, db, task.ID.String())
		if updated.ProjectID != project.ID {
			t.Errorf("expected project ID to be set, %q != %q", updated.ProjectID, project.ID)
		}
	})
}

func MustCreateTask(tb testing.TB, ctx context.Context, db *postgres.DB, cmd domain.CreateTaskCmd) domain.Task {
	tb.Helper()
	task, err := postgres.NewTaskService(db, logrus.New()).Create(ctx, cmd)
	if err != nil {
		tb.Fatal(err)
	}
	return task
}
