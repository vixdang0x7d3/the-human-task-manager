package domain

import "github.com/google/uuid"

type Project struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Title  string
}

type CreateProjectCmd struct {
	Title string
}

type ProjectService interface {
}
