package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/types"
)

type UserHandler struct {
	Service UserService
}

func NewHandler(service UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

type UserService interface {
	CreateUser(ctx context.Context, arg types.CreateUserCmd) (types.User, error)
	ByID(ctx context.Context, id uuid.UUID) (types.User, error)
	ByEmail(ctx context.Context, email string) (types.User, error)
}
