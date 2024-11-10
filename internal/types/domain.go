package types

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserCmd struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string
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
