package core

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type UserCore struct {
	Store UserStore
}

func NewUserCore(store UserStore) *UserCore {
	return &UserCore{
		Store: store,
	}
}

func (c *UserCore) CreateUser(ctx context.Context, cmd types.CreateUserCmd) (types.User, error) {

	hashedPassword, err := hashPassword(cmd.Password)
	if err != nil {
		return types.User{}, err
	}

	user, err := c.Store.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Username:  cmd.Username,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Email:     cmd.Email,
		Password:  hashedPassword,
		SignupAt:  time.Now(),
		LastLogin: time.Now(),
	})
	if err != nil {
		return types.User{}, err
	}

	return toDomainUser(user), nil
}

func (c *UserCore) ByID(ctx context.Context, id uuid.UUID) (types.User, error) {
	user, err := c.Store.ByID(ctx, id)
	if err != nil {
		return types.User{}, err
	}
	return toDomainUser(user), nil
}

func (c *UserCore) ByEmail(ctx context.Context, email string) (types.User, error) {
	user, err := c.Store.ByEmail(ctx, email)
	if err != nil {
		return types.User{}, err
	}
	return toDomainUser(user), nil
}

func (c *UserCore) CheckPassword(ctx context.Context, email string, password string) error {
	u, err := c.Store.ByEmail(ctx, email)
	if err != nil {
		return err
	}
	if !checkPassword(password, u.Password) {
		return errors.New("incorrect password error")
	}
	return nil
}

func toDomainUser(user database.User) types.User {
	return types.User{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		SignupAt:  user.SignupAt,
		LastLogin: user.LastLogin,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type UserStore interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	ByID(ctx context.Context, id uuid.UUID) (database.User, error)
	ByEmail(ctx context.Context, email string) (database.User, error)
}
