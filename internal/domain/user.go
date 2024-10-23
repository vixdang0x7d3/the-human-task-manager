package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
)

// represent core business logic for user domain
type UserCore struct {
	Store UserStore
}

// TODO:
// - field validation especially email, only allows certain characters for password
// - check if email exists already if it is u're cooked
// - check for NOT NULL fields
// - encrypt password
// At the moment the method will just encrypt the password
// and save the user to db
func (core *UserCore) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(arg)

	// this is garbage error handling, however i haven't came up with anything better,
	// so this is what we're having for now
	if err != nil {
		fieldNames := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			fieldNames = append(fieldNames, strings.ToLower(err.Field()))
		}
		return User{}, errors.New("User validation error: Invalid " + strings.Join(fieldNames, ", "))
	}

	hashedPassword, err := internal.HashPassword(arg.Password)
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

// TODO:
// - return a user, check if id exist for fuck sake
// should it accepts a uuid.UUID or a stupid string ?
func (core *UserCore) GetUserByID(userID string) (User, error) {
	return User{}, nil
}

// contract for database layer
type UserStore interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (database.User, error)
}
