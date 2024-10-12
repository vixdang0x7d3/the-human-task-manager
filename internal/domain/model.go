package domain

import (
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
