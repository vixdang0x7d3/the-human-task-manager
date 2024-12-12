package domain

import (
	"context"

	"github.com/google/uuid"
)

type ProjectMembership struct {
	UserID    uuid.UUID
	ProjectID uuid.UUID
	Role      string
}

type ProjectMembershipItem struct {
	User    User
	Project Project
	Role    string
}

type ProjectMembershipService interface {
	Request(ctx context.Context, cmd ProjectMembershipCmd) (ProjectMembership, error)
	Invite(ctx context.Context, cmd ProjectMembershipCmd) (ProjectMembership, error)
	AcceptRequest(ctx context.Context, cmd ProjectMembershipCmd) (ProjectMembership, error)
	AcceptInvitation(ctx context.Context, cmd ProjectMembershipCmd) (ProjectMembership, error)
	Delete(ctx context.Context, cmd ProjectMembershipCmd) (ProjectMembership, error)
	Find(ctx context.Context, filter ProjectMembershipFilter) ([]ProjectMembershipItem, int, error)
}

type ProjectMembershipCmd struct {
	UserID    string
	ProjectID string
	Role      string
}

type ProjectMembershipFilter struct {
	ProjectID *string

	Limit  int
	Offset int
}
