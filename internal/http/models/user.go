package models

type UserView struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	ID        string
}

type TaskViewModel struct {
	Title       string
	Description string
	Create_at   string
	Update_at   string
	Deadline    string
	Schedule    string
	Tags        []string
}
