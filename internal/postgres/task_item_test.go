package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres"
)

// TODO: unimplemented test cases:
// ID not found
// unauthorized access
func TestTaskItemByID(t *testing.T) {

	db := MustOpenDB(t, context.Background())
	s := postgres.NewTaskItemService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME0",
			FirstName: "FIRSTNAME0",
			LastName:  "LASTNAME0",
			Email:     "EMAIL0@email.com",
			Password:  "fuckthefeds123",
		})

		ctxWithUser := domain.NewContextWithUser(context.Background(), &user)

		task := MustCreateTask(t, ctxWithUser, db, domain.CreateTaskCmd{
			ProjectID:   "",
			Description: "TESTTASK0",
			Deadline:    "2024-12-13T07:00",
			Schedule:    "2024-12-13T13:00",
			Wait:        "",
			Priority:    "H",
			Tags:        []string{"school", "INT10187"},
		})

		taskItem, err := s.ByID(ctxWithUser, task.ID.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(task.ID, taskItem.ID) {
			t.Errorf("ID mismatch: %q != %q", task.ID, taskItem.ID)
		}

		if !reflect.DeepEqual(user.ID, taskItem.UserID) {
			t.Errorf("UserID mismatch: %q != %q", user.ID, taskItem.UserID)
		}

		if !reflect.DeepEqual(user.Username, taskItem.Username) {
			t.Errorf("Username mismatch: %q != %q", user.Username, taskItem.Username)
		}

		if taskItem.Deadline.Compare(task.Deadline) != 0 {
			t.Errorf("Deadline mismatch: %v != %v", task.Deadline, taskItem.Deadline)
		}

		if taskItem.Schedule.Compare(task.Schedule) != 0 {
			t.Errorf("Schedule mismatch: %v != %v", task.Schedule, taskItem.Schedule)
		}

		if taskItem.Wait.Compare(task.Wait) != 0 {
			t.Errorf("Wait mismatch: %v != %v", task.Wait, taskItem.Wait)
		}

		if taskItem.Priority != task.Priority {
			t.Errorf("Priority mismatch: %q != %q", task.Priority, taskItem.Priority)
		}

		if !reflect.DeepEqual(taskItem.Tags, task.Tags) {
			t.Errorf("Tags mismatch: %#v != %#v", task.Tags, taskItem.Tags)
		}

		if taskItem.Urgency == 0 {
			t.Errorf("Expect urgency")
		}

	})
}

// TODO: unimplemented test cases:
// invalid filters
// days and months <= 0
// unauthorized access
// pagination ok
func TestTaskItemFind(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	s := postgres.NewTaskItemService(db, logrus.New())

	t.Run("FindByUserID OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME1",
			FirstName: "FIRSTNAME1",
			LastName:  "LASTNAME1",
			Email:     "EMAIL1@email.com",
			Password:  "PASSWORD1",
		})

		ctxWithUser := domain.NewContextWithUser(context.Background(), &user)

		MustCreateTask(t, ctxWithUser, db, domain.CreateTaskCmd{
			ProjectID:   "",
			Description: "TESTTASK1",
			Deadline:    "2024-12-13T07:00",
			Schedule:    "2024-12-13T07:00",
			Wait:        "",
			Priority:    "M",
			Tags:        []string{"school", "INT10187"},
		})
		MustCreateTask(t, ctxWithUser, db, domain.CreateTaskCmd{
			ProjectID:   "",
			Description: "TESTTASK2",
			Deadline:    "2024-12-13T07:00",
			Schedule:    "2024-12-13T07:00",
			Wait:        "",
			Priority:    "M",
			Tags:        []string{"school", "INT10111"},
		})
		MustCreateTask(t, ctxWithUser, db, domain.CreateTaskCmd{
			ProjectID:   "",
			Description: "TESTTASK3",
			Deadline:    "2024-12-13T07:00",
			Schedule:    "2024-12-13T13:00",
			Wait:        "",
			Priority:    "L",
			Tags:        []string{"school", "INT10112"},
		})

		taskItems, n, err := s.Find(
			ctxWithUser,
			domain.TaskItemFilter{
				Limit:  10,
				Offset: 0,
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if n != 3 {
			t.Errorf("mismatch task number returned: %d != %d", 3, len(taskItems))
		}
	})
}

func MustTaskItemByID(tb testing.TB, ctx context.Context, db *postgres.DB, id string) domain.TaskItem {
	tb.Helper()
	taskItem, err := postgres.NewTaskItemService(db, logrus.New()).ByID(ctx, id)
	if err != nil {
		tb.Fatal(err)
	}
	return taskItem
}
