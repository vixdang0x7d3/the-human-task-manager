package types

type UserViewModel struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
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
