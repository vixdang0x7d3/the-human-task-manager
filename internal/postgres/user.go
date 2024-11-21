package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
)

var _ domain.UserService = (*UserService)(nil)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(ctx context.Context, cmd domain.CreateUserCmd) (domain.User, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.User{}, errors.New("cannot get a connection")
	}
	q := sqlc.New(conn)

	user, err := createUser(ctx, q, cmd)
	if err != nil {
		return domain.User{}, err
	}

	return toDomainUser(user), nil
}

func (s *UserService) ByEmail(ctx context.Context, email string) (domain.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.User{}, errors.New("cannot get a connection")
	}
	q := sqlc.New(conn)

	user, err := byEmail(ctx, q, email)
	if err != nil {
		return domain.User{}, err
	}

	return toDomainUser(user), nil
}

func (s *UserService) ByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.User{}, errors.New("cannot get a connection")
	}
	q := sqlc.New(conn)

	user, err := byID(ctx, q, id)
	if err != nil {
		return domain.User{}, err
	}

	return toDomainUser(user), nil
}

func createUser(ctx context.Context, q UserQueries, cmd domain.CreateUserCmd) (sqlc.User, error) {
	panic("unimplemented")
}

func byID(ctx context.Context, q UserQueries, id uuid.UUID) (sqlc.User, error) {
	panic("unimplemented")
}

func byEmail(ctx context.Context, q UserQueries, email string) (sqlc.User, error) {
	panic("unimplemented")
}

func toDomainUser(user sqlc.User) domain.User {
	return domain.User{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		SignupAt:  user.SignupAt,
		LastLogin: user.LastLogin,
	}
}

type UserQueries interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	ByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)
	ByEmail(ctx context.Context, email string) (sqlc.User, error)
}
