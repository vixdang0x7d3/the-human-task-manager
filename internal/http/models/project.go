package models

type ProjectView struct {
	Title  string
	ID     string
	UserID string
}

type ProjectMembershipItemView struct {
	ProjectID string
	UserID    string
	Title     string
	Role      string
	Username  string
}
