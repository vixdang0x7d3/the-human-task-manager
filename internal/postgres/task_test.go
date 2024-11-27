package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
)

func TestCreateTask(t *testing.T) {

	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
		Username:  "chineseman",
		FirstName: "Franz",
		LastName:  "Kafka",
		Email:     "kafka.franz@email.com",
	})

	t.Run("OK", func(t *testing.T) {

		s := postgres.NewTaskService(db)

		cmd := domain.CreateTaskCmd{
			UserID:      user.ID.String(),
			ProjectID:   "",
			Description: "DESCRIPTION",
			Deadline:    "2024-12-23T07:00",
			Schedule:    "",
			Wait:        "",
			Priority:    "",
			Tags:        []string{"TAG1", "TAG2", "TAG3"},
		}

		task, err := s.CreateTask(context.Background(), cmd)
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

func MustCreateTask(tb testing.TB, ctx context.Context, db *postgres.DB, cmd domain.CreateTaskCmd) domain.Task {
	tb.Helper()
	task, err := postgres.NewTaskService(db).CreateTask(ctx, cmd)
	if err != nil {
		tb.Fatal(err)
	}
	return task
}
