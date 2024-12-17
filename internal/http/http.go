package http

import (
	"encoding/json"
	"strconv"

	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/models"
)

func unmarshalTagsJSON(data []byte) (tags []string, err error) {
	if len(data) == 0 {
		return
	}

	var items []struct {
		Value string `json:"value"`
	}
	if err = json.Unmarshal(data, &items); err != nil {
		return []string{}, err
	}

	for _, item := range items {
		tags = append(tags, item.Value)
	}
	return
}

func toTaskItemView(taskItem domain.TaskItem) models.TaskItemView {

	const timeFormat = "2006-01-02T15:04"

	var (
		completedByName string = "none"
		projectTitle    string = "none"
		priority        string = "none"
		deadline        string = "none"
		schedule        string = "none"
		wait            string = "none"
		end             string = "none"
	)

	if taskItem.CompletedByName != "" {
		completedByName = taskItem.CompletedByName
	}

	if taskItem.ProjectTitle != "" {
		projectTitle = taskItem.ProjectTitle
	}

	if taskItem.Priority != "" {
		priority = taskItem.Priority
	}

	if !taskItem.Deadline.IsZero() {
		deadline = taskItem.Deadline.Format(timeFormat)
	}

	if !taskItem.Schedule.IsZero() {
		schedule = taskItem.Schedule.Format(timeFormat)
	}

	if !taskItem.Wait.IsZero() {
		deadline = taskItem.Wait.Format(timeFormat)
	}

	if !taskItem.End.IsZero() {
		deadline = taskItem.End.Format(timeFormat)
	}

	return models.TaskItemView{
		ID:             taskItem.ID.String(),
		Description:    taskItem.Description,
		UserID:         taskItem.UserID.String(),
		Username:       taskItem.Username,
		CompleteBy:     taskItem.CompletedBy.String(),
		CompleteByName: completedByName,
		ProjectID:      taskItem.ProjectID.String(),
		ProjectTitle:   projectTitle,

		Priority: priority,
		State:    taskItem.State,

		Deadline: deadline,
		Schedule: schedule,
		Wait:     wait,
		Create:   taskItem.Create.Format(timeFormat),
		End:      end,
		Urgency:  strconv.FormatFloat(taskItem.Urgency, 'f', 2, 64),

		Tags: taskItem.Tags,
	}
}

func toProjectView(project domain.Project) models.ProjectView {
	return models.ProjectView{
		Title:  project.Title,
		ID:     project.ID.String(),
		UserID: project.UserID.String(),
	}
}
