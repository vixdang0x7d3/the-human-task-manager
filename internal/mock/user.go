package mock

import (
	"context"

	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
)

var _ domain.UserService = (*UserService)(nil)

type UserService struct {
	fnUpdate              func(ctx context.Context, cmd domain.UpdateUserCmd) (domain.User, error)
	fnDelete              func(ctx context.Context) (domain.User, error)
	fnWithPassword        func(ctx context.Context, password string) (domain.User, error)
	fnCreate              func(ctx context.Context, cmd domain.CreateUserCmd) (domain.User, error)
	fnByID                func(ctx context.Context, id string) (domain.User, error)
	fnByEmail             func(ctx context.Context, email string) (domain.User, error)
	fnByEmailWithPassword func(ctx context.Context, email string, password string) (domain.User, error)
}

func (s *UserService) Create(ctx context.Context, cmd domain.CreateUserCmd) (domain.User, error) {
	return s.fnCreate(ctx, cmd)
}

func (s *UserService) ByID(ctx context.Context, id string) (domain.User, error) {
	return s.fnByID(ctx, id)
}

func (s *UserService) ByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.fnByEmail(ctx, email)
}

func (s *UserService) ByEmailWithPassword(ctx context.Context, email string, password string) (domain.User, error) {
	return s.fnByEmailWithPassword(ctx, email, password)
}

func (s *UserService) Update(ctx context.Context, cmd domain.UpdateUserCmd) (domain.User, error) {
	return s.fnUpdate(ctx, cmd)
}

func (s *UserService) Delete(ctx context.Context) (domain.User, error) {
	return s.fnDelete(ctx)
}

func (s *UserService) WithPassword(ctx context.Context, password string) (domain.User, error) {
	return s.fnWithPassword(ctx, password)
}
