package domain

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserParams struct {
	Username  string `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required"`
}

type User struct {
	ID        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	Email     string
	SignupAt  time.Time
	LastLogin time.Time
}
