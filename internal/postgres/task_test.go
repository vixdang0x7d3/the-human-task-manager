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

// TODO: unimplemented test cases:
// user ID not found
// project id ok
// project id not found
// invalid inputs
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
			Description: "DESCRIPTION0",
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

// TODO: unimplemented test cases:
// ID not found
// unauthorized update
// invalid deadline & schedule timestamps
func TestUpdateTask(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	s := postgres.NewTaskService(db)

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME2",
			FirstName: "FIRSTNAME2",
			LastName:  "LASTNAME2",
			Email:     "EMAIL2@email.com",
			Password:  "fckthisshit",
		})

		task := MustCreateTask(t, context.Background(), db, domain.CreateTaskCmd{
			UserID:      user.ID.String(),
			ProjectID:   "",
			Description: "DESCRIPTION1",
			Deadline:    "2024-12-05T07:00",
			Schedule:    "2024-12-02T07:00",
			Wait:        "",
			Priority:    "M",
			Tags:        []string{"school", "INT10111"},
		})

		s.UpdateTask(context.Background(), task.ID.String(), domain.UpdateTaskCmd{
			Description: "ModifiedDescription",
			Deadline:    "",
			Schedule:    "",
			Wait:        "",
			Priority:    "H",
			Tags:        []string{"school"},
		})

		taskItem := MustTaskItemByID(t, context.Background(), db, task.ID.String())

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

	s := postgres.NewTaskService(db)

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME4",
			FirstName: "FIRSTNAME4",
			LastName:  "LASTNAME4",
			Email:     "EMAIL4@email.com",
			Password:  "llalalalala",
		})

		task := MustCreateTask(t, context.Background(), db, domain.CreateTaskCmd{
			UserID:      user.ID.String(),
			ProjectID:   "",
			Description: "DESCRIPTION4",
			Deadline:    "2024-12-05T07:00",
			Schedule:    "2024-12-02T07:00",
			Wait:        "",
			Priority:    "L",
			Tags:        []string{"A", "B"},
		})

		_, err := s.CompleteTask(domain.NewContextWithUser(context.Background(), &user), task.ID.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		completed := MustTaskItemByID(t, context.Background(), db, task.ID.String())
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

	s := postgres.NewTaskService(db)

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME5",
			FirstName: "FIRSTNAME5",
			LastName:  "LASTNAME5",
			Email:     "EMAIL5@email.com",
			Password:  "bababababa",
		})

		task := MustCreateTask(t, context.Background(), db, domain.CreateTaskCmd{
			UserID:      user.ID.String(),
			ProjectID:   "",
			Description: "DESCRIPTION5",
			Deadline:    "2024-12-05T07:00",
			Schedule:    "2024-12-02T07:00",
			Wait:        "",
			Priority:    "L",
			Tags:        []string{"B", "C"},
		})

		_, err := s.DeleteTask(context.Background(), task.ID.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		deleted := MustTaskItemByID(t, context.Background(), db, task.ID.String())
		if deleted.State != "deleted" {
			t.Errorf("wrong state: %q != %q", deleted.State, "deleted")
		}
	})
}

func TestStartWaitingTask(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewTaskService(db)

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME6",
			FirstName: "FIRSTNAME6",
			LastName:  "LASTNAME6",
			Email:     "EMAIL6@email.com",
			Password:  "yayayayay",
		})

		startMe1 := MustCreateTask(t, context.Background(), db, domain.CreateTaskCmd{
			UserID:      user.ID.String(),
			ProjectID:   "",
			Description: "STARTME1",
			Deadline:    "",
			Schedule:    "",
			Wait:        "2024-11-30T11:00",
			Priority:    "H",
			Tags:        []string{},
		})

		startMe2 := MustCreateTask(t, context.Background(), db, domain.CreateTaskCmd{
			UserID:      user.ID.String(),
			ProjectID:   "",
			Description: "STARTME2",
			Deadline:    "",
			Schedule:    "",
			Wait:        "2024-11-27T07:00",
			Priority:    "H",
			Tags:        []string{},
		})

		dontStart := MustCreateTask(t, context.Background(), db, domain.CreateTaskCmd{
			UserID:      user.ID.String(),
			ProjectID:   "",
			Description: "DONTTOUCHME",
			Deadline:    "",
			Schedule:    "",
			Wait:        "2025-01-01T07:00",
			Priority:    "M",
			Tags:        []string{},
		})

		started, err := s.StartWaitingTasks(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(started) != 2 {
			t.Errorf("expected 2 tasks to start, got %d", len(started))
		}

		if MustTaskItemByID(t, context.Background(), db, startMe1.ID.String()).State != "started" {
			t.Errorf("expected %s get started", startMe1.Description)
		}

		if MustTaskItemByID(t, context.Background(), db, startMe2.ID.String()).State != "started" {
			t.Errorf("expected %s get started", startMe2.Description)
		}

		if MustTaskItemByID(t, context.Background(), db, dontStart.ID.String()).State != "waiting" {
			t.Errorf("expected %s stay waiting", dontStart.Description)
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
