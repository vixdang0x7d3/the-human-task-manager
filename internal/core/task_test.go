package core

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/types"
)

type StubTaskStore struct {
	Tasks map[uuid.UUID]database.Task
}

func (s *StubTaskStore) CreateTask(ctx context.Context, arg database.CreateTaskParams) (database.Task, error) {

	// mocking database behaviour when fetching null uuid
	valid := !reflect.DeepEqual(arg.ProjectID.UUID, uuid.Nil)
	return database.Task{
		ID:     arg.ID,
		UserID: arg.UserID,
		ProjectID: uuid.NullUUID{
			UUID:  arg.ProjectID.UUID,
			Valid: valid,
		},
		Description: arg.Description,
		Priority:    arg.Priority,
		Status:      arg.Status,
		Deadline:    arg.Deadline,
		Schedule:    arg.Schedule,
		Wait:        arg.Wait,
		Create:      arg.Create,
		End:         arg.End,
	}, nil
}

func TestCreateTask(t *testing.T) {

	mockstore := &StubTaskStore{
		map[uuid.UUID]database.Task{},
	}

	t.Run("it saves new task to store and return the saved data as domain task", func(t *testing.T) {
		arg := types.CreateTaskCmd{
			UserID:      "7503f390-81b9-4ac4-8036-88cc52901380",
			ProjectID:   "",
			Description: "a task used for testing",
			Deadline:    "2024-11-25T10:00",
			Schedule:    "",
			Wait:        "",
			Priority:    "",
		}
		_ = arg

		core := NewTaskCore(mockstore)
		task, err := core.CreateTask(context.Background(), arg)
		if err != nil {
			t.Fatalf(`Unexpected error: "%v"`, err)
		}

		wantDeadline, _ := time.Parse(time.RFC3339, "2024-11-25T10:00:00Z")
		wantUserID := uuid.MustParse("7502f390-81b9-4ac4-8036-88cc52901380")

		if wantDeadline.Compare(task.Deadline) != 0 {
			t.Errorf(`failed to parse deadline correctly, want "%v", got "%v"`, wantDeadline, task.Deadline)
		}

		if reflect.DeepEqual(wantUserID, task.UserID) {
			t.Errorf(`failed to parse user id correctly, want "%v", got "%v"`, wantUserID, task.UserID)
		}

		if !reflect.DeepEqual(task.ProjectID, uuid.Nil) {
			t.Errorf(`default project id should be uuid.Nil, want "%v", got "%v"`, uuid.Nil, task.ProjectID)
		}

		if task.Priority != "" {
			t.Errorf(`default priority should be zero value, want "%s" got "%s"`, "", task.Priority)
		}

		if task.Status != string(database.TaskStatusStarted) {
			t.Errorf(`default task status should be "started", want "%s", got "%s"`, string(database.TaskStatusStarted), task.Status)
		}

		if !task.Schedule.IsZero() {
			t.Errorf(`default task schedule should be zero value, want "%v", got "%v"`, time.Time{}, task.Schedule)
		}

		if !task.Wait.IsZero() {
			t.Errorf(`default task wait should be zero value, want "%v", got "%v"`, time.Time{}, task.Wait)
		}

	})

	t.Run("it returns a non nil project id when input has project id", func(t *testing.T) {
		arg := types.CreateTaskCmd{
			UserID:      "7503f390-81b9-4ac4-8036-88cc52901380",
			ProjectID:   "defdc269-2046-49f9-87bc-766218588728",
			Description: "a task used for testing project id",
			Deadline:    "",
			Schedule:    "",
			Wait:        "2024-11-25T10:00",
			Priority:    "L",
		}
		_ = arg

		core := NewTaskCore(mockstore)
		task, err := core.CreateTask(context.Background(), arg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := uuid.MustParse("defdc269-2046-49f9-87bc-766218588728")
		if !reflect.DeepEqual(task.ProjectID, want) {
			t.Errorf(`expected returning valid project id, want "%v" got "%v"`, want, task.ProjectID)
		}
	})

	t.Run("it returns error when policies are violated", func(t *testing.T) {
		tests := []struct {
			Description string
			Input       types.CreateTaskCmd
			ErrorMsg    string
		}{
			{
				Description: "error parsing user id",
				Input: types.CreateTaskCmd{
					UserID:      "invaliduserid",
					Description: "",
					Deadline:    "",
					Schedule:    "",
					Wait:        "",
					Priority:    "",
				},
				ErrorMsg: "failed to parse userID",
			},

			{
				Description: "error parsing deadline",
				Input: types.CreateTaskCmd{
					UserID:      "7502f390-81b9-4ac4-8036-88cc52901380",
					Description: "task with errornous deadline",
					Deadline:    "2024/11/123",
					Schedule:    "",
					Wait:        "",
					Priority:    "",
				},
				ErrorMsg: "failed to parse deadline",
			},

			{
				Description: "error when deadline before time.now",
				Input: types.CreateTaskCmd{
					UserID:      "7502f390-81b9-4ac4-8036-88cc52901380",
					Description: "task with errornous deadline",
					Deadline:    "2023-01-10T07:30",
					Schedule:    "",
					Wait:        "",
					Priority:    "",
				},
				ErrorMsg: "deadline with invalid timestamp",
			},

			{
				Description: "error parsing priority",
				Input: types.CreateTaskCmd{
					UserID:      "7502f390-81b9-4ac4-8036-88cc52901380",
					Description: "task with invalid priority as input",
					Deadline:    "2023-11-25T07:30",
					Schedule:    "",
					Wait:        "",
					Priority:    "invalidpriority",
				},
				ErrorMsg: "failed to parse priority",
			},
		}

		for _, test := range tests {
			t.Run(test.Description, func(t *testing.T) {

				core := NewTaskCore(mockstore)
				_, err := core.CreateTask(context.Background(), test.Input)
				if err == nil {
					t.Errorf(`should return an non-nil error, want err.Error() = "%v" got %v`, test.ErrorMsg, err)
				}
			})
		}
	})
}
