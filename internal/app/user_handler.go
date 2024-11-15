package app

import (
	"context"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/types"
)

type UserHandler struct {
	Service        UserService
	SessionManager *scs.SessionManager
}

func NewUserHandler(service UserService, sessionManager *scs.SessionManager) *UserHandler {
	return &UserHandler{
		Service:        service,
		SessionManager: sessionManager,
	}
}

type UserService interface {
	CreateUser(ctx context.Context, arg types.CreateUserCmd) (types.User, error)
	ByID(ctx context.Context, id uuid.UUID) (types.User, error)
	ByEmail(ctx context.Context, email string) (types.User, error)
	CheckPassword(ctx context.Context, email string, password string) (types.User, error)
}
