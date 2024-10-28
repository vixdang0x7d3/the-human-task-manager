package app

import (
	"time"

	"github.com/google/uuid"
)

// explicit AppUser type used for rendering view
type AppUser struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	SignupAt  time.Time `json:"signup_at"`
	LastLogin time.Time `json:"last_login"`
}
