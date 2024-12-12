package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	Email     string
	SignupAt  time.Time
	LastLogin time.Time
}

type UserService interface {
	Create(ctx context.Context, cmd CreateUserCmd) (User, error)
	ByID(ctx context.Context, id string) (User, error)
	ByEmail(ctx context.Context, email string) (User, error)
	ByEmailWithPassword(ctx context.Context, email string, password string) (User, error)
}

type CreateUserCmd struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string
}
