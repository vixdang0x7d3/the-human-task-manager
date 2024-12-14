package models

type TaskView struct {
	ID          string
	Description string
	UserID      string
	Username    string
	ProjectID   string
	ProjectName string

	Priority string
	State    string

	Deadline string
	Schedule string
	Wait     string
	Create   string
	End      string
	Tags     []string
}

type TaskItemView struct {
	ID             string
	Description    string
	Username       string
	UserID         string
	CompleteBy     string
	CompleteByName string
	ProjectTitle   string
	ProjectID      string

	Priority string
	State    string

	Deadline string
	Schedule string
	Wait     string
	Create   string
	End      string

	Tags    []string
	Urgency string
}
