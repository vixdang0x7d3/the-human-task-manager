package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres/sqlc"
)

var _ domain.ProjectMembershipService = (*ProjectMembershipService)(nil)

type ProjectMembershipService struct {
	db     *DB
	logger *logrus.Logger
}

func NewProjectMembershipService(db *DB, logger *logrus.Logger) *ProjectMembershipService {
	return &ProjectMembershipService{
		db:     db,
		logger: logger,
	}
}

func (s *ProjectMembershipService) AcceptInvitation(
	ctx context.Context,
	cmd domain.ProjectMembershipCmd,
) (domain.ProjectMembership, error) {

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return domain.ProjectMembership{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
	}
	cmd.UserID = (*userID).String()

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.ProjectMembership{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	membership, err := acceptMembershipInvitation(ctx, q, cmd)
	if err != nil {
		return domain.ProjectMembership{}, err
	}

	return toDomainMembership(membership), nil
}

func (s *ProjectMembershipService) AcceptRequest(
	ctx context.Context,
	cmd domain.ProjectMembershipCmd,
) (domain.ProjectMembership, error) {

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return domain.ProjectMembership{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "No user ID in context",
		}
	}

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.ProjectMembership{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)
	membership, err := acceptMembershipRequest(ctx, q, cmd)
	if err != nil {
		return domain.ProjectMembership{}, err
	}

	return toDomainMembership(membership), nil
}

func (s *ProjectMembershipService) Invite(
	ctx context.Context,
	cmd domain.ProjectMembershipCmd,
) (domain.ProjectMembership, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.ProjectMembership{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return domain.ProjectMembership{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
	}

	project, err := projectByID(ctx, q, cmd.ProjectID)
	if err != nil {
		return domain.ProjectMembership{}, err
	}

	if project.UserID != *userID {
		return domain.ProjectMembership{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "unauthorized invite, not project owner",
		}
	}

	membership, err := inviteMembership(ctx, q, cmd)
	if err != nil {
		return domain.ProjectMembership{}, err
	}

	return toDomainMembership(membership), nil
}

func (s *ProjectMembershipService) Request(
	ctx context.Context,
	cmd domain.ProjectMembershipCmd,
) (domain.ProjectMembership, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.ProjectMembership{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return domain.ProjectMembership{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
	}
	cmd.UserID = (*userID).String()

	project, err := projectByID(ctx, q, cmd.ProjectID)
	if err != nil {
		return domain.ProjectMembership{}, err
	}

	if project.UserID == *userID {
		return domain.ProjectMembership{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "cannot request membership, is owner",
		}
	}

	membership, err := requestMembership(ctx, q, cmd)
	if err != nil {
		return domain.ProjectMembership{}, err
	}

	return toDomainMembership(membership), nil
}

func (s *ProjectMembershipService) Delete(
	ctx context.Context,
	cmd domain.ProjectMembershipCmd,
) (domain.ProjectMembership, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return domain.ProjectMembership{}, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return domain.ProjectMembership{}, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
	}

	project, err := projectByID(ctx, q, cmd.ProjectID)
	if err != nil {
		return domain.ProjectMembership{}, err
	}

	// when current user is the project owner,
	// they can not delete themself
	if *userID == project.UserID {
		if memberID, err := uuid.Parse(cmd.UserID); err != nil {
			return domain.ProjectMembership{}, &domain.Error{
				Code:    domain.EINVALID,
				Message: "corrupted user ID",
			}
		} else if *userID == memberID {
			return domain.ProjectMembership{}, &domain.Error{
				Code:    domain.ECONFLICT,
				Message: "cannot delete membership, user is project owner",
			}
		}

		membership, err := deleteMembership(ctx, q, domain.ProjectMembershipCmd{
			UserID:    cmd.UserID,
			ProjectID: cmd.ProjectID,
		})
		if err != nil {
			return domain.ProjectMembership{}, err
		}

		return toDomainMembership(membership), nil
	}

	// when current user is not the project owner,
	// they can only delete themself
	// in this case, the userID in cmd should be empty
	// {
	//	UserID: "",
	//	ProjectID: projectID
	// }
	membership, err := deleteMembership(ctx, q, domain.ProjectMembershipCmd{
		UserID:    (*userID).String(),
		ProjectID: cmd.ProjectID,
	})
	if err != nil {
		return domain.ProjectMembership{}, err
	}

	return toDomainMembership(membership), nil
}

// Find returns a list of embedded memberships that matches the filter.
// If filter.ProjectID is not set, Find will return a list of memberships
// of the current user.
// If filter.ProjectID is set, Find will return memberships of the project
// but if only the current user is a member of that project.
func (s *ProjectMembershipService) Find(
	ctx context.Context,
	filter domain.ProjectMembershipFilter,
) ([]domain.ProjectMembershipItem, int, error) {

	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return []domain.ProjectMembershipItem{}, 0, err
	}
	defer conn.Release()

	q := sqlc.New(conn)

	if filter.ProjectID != nil {
		rows, n, err := findMembershipsByProjectID(ctx, q, filter)
		if err != nil {
			return []domain.ProjectMembershipItem{}, 0, err
		}

		return Map(rows, byProjectIDRow_toDomainMembershipItem), int(n), nil
	}

	rows, n, err := findMembershipsByUserID(ctx, q, filter)
	if err != nil {
		return []domain.ProjectMembershipItem{}, 0, err
	}

	return Map(rows, byUserIDRow_toDomainMembershipItem), int(n), nil
}

func acceptMembershipInvitation(
	ctx context.Context,
	q ProjectMembershipQueries,
	cmd domain.ProjectMembershipCmd,
) (sqlc.ProjectMembership, error) {

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted user ID",
		}
	}

	projectID, err := uuid.Parse(cmd.ProjectID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted project ID",
		}
	}

	membership, err := q.AcceptMembership(ctx, sqlc.AcceptMembershipParams{
		UserID:    userID,
		ProjectID: projectID,
	})
	if err != nil {
		return sqlc.ProjectMembership{}, err
	}

	return membership, nil
}

func acceptMembershipRequest(
	ctx context.Context,
	q ProjectMembershipQueries,
	cmd domain.ProjectMembershipCmd,
) (sqlc.ProjectMembership, error) {

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted user ID",
		}
	}

	projectID, err := uuid.Parse(cmd.ProjectID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted project ID",
		}
	}

	membership, err := q.AcceptMembership(ctx, sqlc.AcceptMembershipParams{
		UserID:    userID,
		ProjectID: projectID,
	})
	if err != nil {
		return sqlc.ProjectMembership{}, err
	}

	return membership, nil
}

func inviteMembership(
	ctx context.Context,
	q ProjectMembershipQueries,
	cmd domain.ProjectMembershipCmd,
) (sqlc.ProjectMembership, error) {

	guestID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted user ID of project guest",
		}
	}

	projectID, err := uuid.Parse(cmd.ProjectID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted project ID",
		}
	}

	membership, err := q.CreateMembership(ctx, sqlc.CreateMembershipParams{
		UserID:    guestID,
		ProjectID: projectID,
		Role:      sqlc.MembershipRoleInvited,
	})
	if err != nil {
		return sqlc.ProjectMembership{}, err
	}

	return membership, nil
}

func requestMembership(
	ctx context.Context,
	q ProjectMembershipQueries,
	cmd domain.ProjectMembershipCmd,
) (sqlc.ProjectMembership, error) {

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted user ID",
		}
	}

	projectID, err := uuid.Parse(cmd.ProjectID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted project ID",
		}
	}

	membership, err := q.CreateMembership(
		ctx,
		sqlc.CreateMembershipParams{
			UserID:    userID,
			ProjectID: projectID,
			Role:      sqlc.MembershipRoleRequested,
		},
	)
	if err != nil {
		return sqlc.ProjectMembership{}, err
	}

	return membership, nil
}

func createOwnerMembership(
	ctx context.Context,
	q ProjectMembershipQueries,
	cmd domain.ProjectMembershipCmd,
) (sqlc.ProjectMembership, error) {

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted user ID",
		}
	}

	projectID, err := uuid.Parse(cmd.ProjectID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted project ID",
		}
	}

	membership, err := q.CreateMembership(ctx, sqlc.CreateMembershipParams{
		UserID:    userID,
		ProjectID: projectID,
		Role:      sqlc.MembershipRoleOwner,
	})
	if err != nil {
		return sqlc.ProjectMembership{}, err
	}

	return membership, nil
}

func deleteMembership(
	ctx context.Context,
	q ProjectMembershipQueries,
	cmd domain.ProjectMembershipCmd,
) (sqlc.ProjectMembership, error) {

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted user ID",
		}
	}

	projectID, err := uuid.Parse(cmd.ProjectID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted project ID",
		}
	}

	membership, err := q.DeleteMembership(ctx, sqlc.DeleteMembershipParams{
		UserID:    userID,
		ProjectID: projectID,
	})
	if err != nil {
		return sqlc.ProjectMembership{}, err
	}

	return membership, nil
}

func findMembershipsByUserID(
	ctx context.Context,
	q ProjectMembershipQueries,
	filter domain.ProjectMembershipFilter,
) ([]sqlc.MembershipsByUserIDRow, int64, error) {

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return []sqlc.MembershipsByUserIDRow{}, 0, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
	}

	rows, err := q.MembershipsByUserID(ctx, sqlc.MembershipsByUserIDParams{
		UserID:  *userID,
		Nlimit:  int32(filter.Limit),
		Noffset: int32(filter.Offset),
	})
	if err != nil {
		return []sqlc.MembershipsByUserIDRow{}, 0, err
	}

	if len(rows) == 0 {
		return []sqlc.MembershipsByUserIDRow{}, 0, &domain.Error{
			Code:    domain.ENOTFOUND,
			Message: "no memberships found",
		}
	}

	return rows, rows[0].Count, nil
}

func findMembershipsByProjectID(
	ctx context.Context,
	q ProjectMembershipQueries,
	filter domain.ProjectMembershipFilter,
) ([]sqlc.MembershipsByProjectIDRow, int64, error) {

	userID := domain.UserIDFromContext(ctx)
	if userID == nil {
		return []sqlc.MembershipsByProjectIDRow{}, 0, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "no user ID in context",
		}
	}

	if membership, err := membershipByIDs(ctx, q, domain.ProjectMembershipCmd{
		UserID:    userID.String(),
		ProjectID: *filter.ProjectID,
	}); err != nil {
		return []sqlc.MembershipsByProjectIDRow{}, 0, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "unauthorized access, not a project member",
		}

	} else if membership.Role == sqlc.MembershipRoleInvited || membership.Role == sqlc.MembershipRoleRequested {
		return []sqlc.MembershipsByProjectIDRow{}, 0, &domain.Error{
			Code:    domain.EUNAUTHORIZED,
			Message: "unauthorized access, not a project member yet",
		}
	}

	projectID, err := uuid.Parse(*filter.ProjectID)
	if err != nil {
		return []sqlc.MembershipsByProjectIDRow{}, 0, err
	}

	rows, err := q.MembershipsByProjectID(ctx, sqlc.MembershipsByProjectIDParams{
		ProjectID: projectID,
		Nlimit:    int32(filter.Limit),
		Noffset:   int32(filter.Offset),
	})
	if err != nil {
		return []sqlc.MembershipsByProjectIDRow{}, 0, err
	}

	if len(rows) == 0 {
		return []sqlc.MembershipsByProjectIDRow{}, 0, &domain.Error{
			Code:    domain.ENOTFOUND,
			Message: "no memberships found",
		}
	}

	return rows, rows[0].Count, nil
}

func membershipByIDs(
	ctx context.Context,
	q ProjectMembershipQueries,
	cmd domain.ProjectMembershipCmd,
) (sqlc.ProjectMembership, error) {

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted user ID",
		}
	}

	projectID, err := uuid.Parse(cmd.ProjectID)
	if err != nil {
		return sqlc.ProjectMembership{}, &domain.Error{
			Code:    domain.EINVALID,
			Message: "corrupted project ID",
		}
	}

	membership, err := q.MembershipByIDs(ctx, sqlc.MembershipByIDsParams{
		UserID:    userID,
		ProjectID: projectID,
	})
	if err != nil {
		return sqlc.ProjectMembership{}, err
	}

	return membership, nil
}

func toDomainMembership(membership sqlc.ProjectMembership) domain.ProjectMembership {
	return domain.ProjectMembership{
		ProjectID: membership.ProjectID,
		UserID:    membership.UserID,
		Role:      string(membership.Role),
	}
}

func byUserIDRow_toDomainMembershipItem(row sqlc.MembershipsByUserIDRow) domain.ProjectMembershipItem {
	return domain.ProjectMembershipItem{
		User:    toDomainUser(row.User),
		Project: toDomainProject(row.Project),
		Role:    string(row.Role),
	}
}

func byProjectIDRow_toDomainMembershipItem(row sqlc.MembershipsByProjectIDRow) domain.ProjectMembershipItem {
	return domain.ProjectMembershipItem{
		User:    toDomainUser(row.User),
		Project: toDomainProject(row.Project),
		Role:    string(row.Role),
	}

}

type ProjectMembershipQueries interface {
	CreateMembership(ctx context.Context, arg sqlc.CreateMembershipParams) (sqlc.ProjectMembership, error)
	AcceptMembership(ctx context.Context, arg sqlc.AcceptMembershipParams) (sqlc.ProjectMembership, error)
	DeleteMembership(ctx context.Context, arg sqlc.DeleteMembershipParams) (sqlc.ProjectMembership, error)
	MembershipByIDs(ctx context.Context, arg sqlc.MembershipByIDsParams) (sqlc.ProjectMembership, error)
	MembershipsByUserID(ctx context.Context, arg sqlc.MembershipsByUserIDParams) ([]sqlc.MembershipsByUserIDRow, error)
	MembershipsByProjectID(ctx context.Context, arg sqlc.MembershipsByProjectIDParams) ([]sqlc.MembershipsByProjectIDRow, error)
}
