package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
	"golang.org/x/crypto/bcrypt"
)

var _ domain.UserService = (*UserService)(nil)

type UserService struct {
	db     *DB
	logger *logrus.Logger
}

func NewUserService(db *DB, logger *logrus.Logger) *UserService {
	return &UserService{
		db:     db,
		logger: logger,
	}
}

func (s *UserService) Create(ctx context.Context, cmd domain.CreateUserCmd) (domain.User, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	user, err := createUser(ctx, q, cmd)
	if err != nil {
		return toDomainUser(user), err
	}

	return toDomainUser(user), nil
}

func (s *UserService) Update(ctx context.Context, cmd domain.UpdateUserCmd) (domain.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	user, err := updateUser(ctx, q, cmd)
	if err != nil {
		return toDomainUser(user), err
	}
	return toDomainUser(user), nil
}

func (s *UserService) ByEmail(ctx context.Context, email string) (domain.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	user, err := userByEmail(ctx, q, email)
	if err != nil {
		return toDomainUser(user), err
	}

	return toDomainUser(user), nil
}

func (s *UserService) ByEmailWithPassword(ctx context.Context, email string, password string) (domain.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	user, err := userByEmailWithPassword(ctx, q, email, password)
	if err != nil {
		return toDomainUser(user), err
	}

	return toDomainUser(user), nil
}

func (s *UserService) ByID(ctx context.Context, id string) (domain.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	user, err := userByID(ctx, q, id)
	if err != nil {
		return toDomainUser(user), err
	}

	return toDomainUser(user), nil
}

func (s *UserService) WithPassword(ctx context.Context, password string) (domain.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	user, err := userWithPassword(ctx, q, password)
	if err != nil {
		return toDomainUser(user), err
	}

	return toDomainUser(user), nil
}

func createUser(ctx context.Context, q UserQueries, cmd domain.CreateUserCmd) (sqlc.User, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, err
	}

	user, err := q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:        uuid.New(),
		Username:  cmd.Username,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Email:     cmd.Email,
		Password:  string(bytes),
		SignupAt:  time.Now(),
		LastLogin: time.Now(),
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return sqlc.User{}, err
		}
		switch pgErr.Code {
		case pgErrCode_UniqueViolation:
			return sqlc.User{}, &domain.Error{
				Code:    domain.ECONFLICT,
				Message: "createUser: this email exists",
			}
		default:
			return sqlc.User{}, err
		}
	}

	return user, nil
}

// updateUser is an authenticated endpoint and
// it checks context for the current logged in user
func updateUser(ctx context.Context, q UserQueries, cmd domain.UpdateUserCmd) (sqlc.User, error) {

	var (
		err   error
		bytes []byte
	)

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return sqlc.User{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
	}
	user, err := userByID(ctx, q, userID.String())

	if cmd.Username != nil {
		user.Username = *cmd.Username
	}

	if cmd.FirstName != nil {
		user.FirstName = *cmd.FirstName
	}

	if cmd.LastName != nil {
		user.LastName = *cmd.LastName
	}

	if cmd.Email != nil {
		user.Email = *cmd.Email
	}

	if cmd.Password != nil {
		bytes, err = bcrypt.GenerateFromPassword([]byte(*cmd.Password), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, nil
		}

		user.Password = string(bytes)
	}

	user, err = q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return sqlc.User{}, err
		}

		switch pgErr.Code {
		case pgErrCode_UniqueViolation:
			return sqlc.User{}, &domain.Error{
				Code:    domain.ECONFLICT,
				Message: "this email exists",
			}
		default:
			return sqlc.User{}, err
		}
	}
	return user, nil
}

func userByID(ctx context.Context, q UserQueries, id string) (sqlc.User, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return sqlc.User{}, &domain.Error{Code: domain.EINVALID, Message: "corrupted ID"}
	}

	user, err := q.UserByID(ctx, uuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.User{}, &domain.Error{
				Code:    domain.ENOTFOUND,
				Message: "userByID: ID not found",
			}
		}
		return sqlc.User{}, err
	}
	return user, nil
}

func userWithPassword(ctx context.Context, q UserQueries, password string) (sqlc.User, error) {
	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return sqlc.User{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "userWithPassword: no user ID in context",
		}
	}

	user, err := userByID(ctx, q, userID.String())
	if err != nil {
		return sqlc.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return sqlc.User{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "userWithPassword: wrong password",
		}
	}

	return user, nil
}

func userByEmail(ctx context.Context, q UserQueries, email string) (sqlc.User, error) {
	user, err := q.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.User{}, &domain.Error{
				Code:    domain.ENOTFOUND,
				Message: "userByEmail: email not found",
			}
		}
		return sqlc.User{}, err
	}

	return user, nil
}

func userByEmailWithPassword(ctx context.Context, q UserQueries, email, password string) (sqlc.User, error) {
	user, err := q.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.User{}, &domain.Error{
				Code:    domain.ENOTFOUND,
				Message: "userByEmailWithPassword: email not found",
			}
		}
		return sqlc.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return sqlc.User{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "userByEmailWithPassword: wrong password",
		}
	}

	return user, nil
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
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	UserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)
	UserByEmail(ctx context.Context, email string) (sqlc.User, error)
}
