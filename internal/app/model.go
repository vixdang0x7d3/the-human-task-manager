package app

import (
	"time"

	"github.com/google/uuid"
)

// AppUser is the View Model for User domain
type AppUser struct {
	ID        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	Email     string
	SignupAt  time.Time
	LastLogin time.Time
}
