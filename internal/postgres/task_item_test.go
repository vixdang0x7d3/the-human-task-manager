package postgres_test

import (
	"context"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres"
	"reflect"
	"testing"
)

func TestTaskItemByID(t *testing.T) {

	db := MustOpenDB(t, context.Background())
	s := postgres.NewTaskItemService(db)

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "kingterry",
			FirstName: "Terry",
			LastName:  "David",
			Email:     "kingterry@coolkid.com",
			Password:  "fuckthefeds123",
		})
		task := MustCreateTask(t, context.Background(), db, domain.CreateTaskCmd{
			UserID:      user.ID.String(),
			ProjectID:   "",
			Description: "build an app",
			Deadline:    "2024-12-03T07:00",
			Schedule:    "2024-11-28T13:00",
			Wait:        "",
			Priority:    "H",
			Tags:        []string{"school", "INT10187"},
		})

		taskItem, err := s.TaskItemByID(context.Background(), task.ID.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(task.ID, taskItem.ID) {
			t.Errorf("ID mismatch: %q != %q", task.ID, taskItem.ID)
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
