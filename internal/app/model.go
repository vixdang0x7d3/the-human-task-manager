package app

import (
	"time"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
)

type AppUser struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	SignupAt  time.Time `json:"signup_at"`
	LastLogin time.Time `json:"last_login"`
}

type UserFormPayload struct {
	Username  string `form:"username"`
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
	Email     string `form:"email"`
	Password  string `form:"password"`
}

func toAppUser(user domain.User) AppUser {
	return AppUser(user)
}
