package models

type TaskView struct {
	Title       string
	Description string
	Tags        []string
	Deadline    string
	Schedule    string
	CreatedAt   string
	UpdatedAt   string
}
