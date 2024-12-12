package domain

import (
	"context"

	"github.com/google/uuid"
)

type contextKey int

const (
	_ = contextKey(iota)
	userContextKey
)

func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) *User {
	user, _ := ctx.Value(userContextKey).(*User)
	return user
}

func UserIDFromContext(ctx context.Context) *uuid.UUID {
	if user := UserFromContext(ctx); user != nil {
		return &user.ID
	}
	return nil
}
