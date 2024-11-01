package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
)

// represent core business logic for user domain
type UserCore struct {
	Store UserStore
}

func (core *UserCore) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {

	hashedPassword, err := HashPassword(arg.Password)
	if err != nil {
		return User{}, err
	}

	user, err := core.Store.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Username:  arg.Username,
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		Email:     arg.Email,
		Password:  hashedPassword,
		SignupAt:  time.Now(),
		LastLogin: time.Now(),
	})
	if err != nil {
		return User{}, err
	}

	return toDomainUser(user), nil
}

// it so simple so im not gonna write test :P
func (core *UserCore) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := core.Store.GetUser(ctx, id)
	if err != nil {
		return User{}, err
	}
	return toDomainUser(user), nil
}

func toDomainUser(user database.User) User {
	return User{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		SignupAt:  user.SignupAt,
		LastLogin: user.LastLogin,
	}
}

type UserStore interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (database.User, error)
}
