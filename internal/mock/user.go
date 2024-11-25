package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
)

var _ domain.UserService = (*UserService)(nil)

type UserService struct {
	fnCreateUser          func(ctx context.Context, cmd domain.CreateUserCmd) (domain.User, error)
	fnByID                func(ctx context.Context, id uuid.UUID) (domain.User, error)
	fnByEmail             func(ctx context.Context, email string) (domain.User, error)
	fnByEmailWithPassword func(ctx context.Context, email string, password string) (domain.User, error)
}

func (s *UserService) CreateUser(ctx context.Context, cmd domain.CreateUserCmd) (domain.User, error) {
	return s.fnCreateUser(ctx, cmd)
}

func (s *UserService) ByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.fnByID(ctx, id)
}

func (s *UserService) ByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.fnByEmail(ctx, email)
}

func (s *UserService) ByEmailWithPassword(ctx context.Context, email string, password string) (domain.User, error) {
	return s.fnByEmailWithPassword(ctx, email, password)
}
