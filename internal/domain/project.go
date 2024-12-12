package domain

import (
	"context"

	"github.com/google/uuid"
)

type Project struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Title  string
}

type ProjectService interface {
	Find(ctx context.Context, filter ProjectFilter) ([]Project, int, error)
	ByID(ctx context.Context, id string) (Project, error)
	Create(ctx context.Context, title string) (Project, error)
	Delete(ctx context.Context, id string) (Project, error)
}

type ProjectFilter struct {
	Limit  int
	Offset int
}
