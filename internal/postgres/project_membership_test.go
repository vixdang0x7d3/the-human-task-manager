package postgres_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres"
)

func TestRequestProjectMembership(t *testing.T) {

	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewProjectMembershipService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		owner := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME12",
			FirstName: "FIRSTNAME12",
			LastName:  "LASTNAME12",
			Email:     "EMAIL12@email.com",
			Password:  "papayaga",
		})

		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME11",
			FirstName: "FIRSTNAME11",
			LastName:  "LASTNAME11",
			Email:     "EMAIL11@email.com",
			Password:  "papayaga",
		})
		project := MustCreateProject(t, domain.NewContextWithUser(context.Background(), &owner), db, "TESTPROJECT")

		projectMembership, err := s.Request(
			domain.NewContextWithUser(context.Background(), &user),
			domain.ProjectMembershipCmd{
				ProjectID: project.ID.String(),
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if projectMembership.UserID != user.ID {
			t.Errorf("user ID mismatches: %q != %q", projectMembership.UserID, user.ID)
		}
	})
}

func TestInviteProjectMembership(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewProjectMembershipService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME14",
			FirstName: "FIRSTNAME14",
			LastName:  "LASTNAME14",
			Email:     "EMAIL14@email.com",
			Password:  "papayaga",
		})
		guest := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME13",
			FirstName: "FIRSTNAME13",
			LastName:  "LASTNAME13",
			Email:     "EMAIL13@email.com",
			Password:  "papayaga",
		})
		project := MustCreateProject(t, domain.NewContextWithUser(context.Background(), &user), db, "TESTPROJECT")

		projectMembership, err := s.Invite(
			domain.NewContextWithUser(context.Background(), &user),
			domain.ProjectMembershipCmd{
				UserID:    guest.ID.String(),
				ProjectID: project.ID.String(),
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if projectMembership.UserID != guest.ID {
			t.Errorf("user ID mismatches: %q != %q", projectMembership.ProjectID, guest.ID)
		}
	})
}

func TestAcceptInvitation(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewProjectMembershipService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME15",
			FirstName: "FIRSTNAME15",
			LastName:  "LASTNAME15",
			Email:     "EMAIL15@email.com",
			Password:  "urmomfat",
		})
		guest := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME16",
			FirstName: "FIRSTNAME16",
			LastName:  "LASTNAME16",
			Email:     "EMAIL16@email.com",
			Password:  "theboysavior",
		})

		project := MustCreateProject(t, domain.NewContextWithUser(context.Background(), &user), db, "TESTPROJECT2")

		projectMembership := MustInvite(
			t,
			domain.NewContextWithUser(context.Background(), &user),
			db,
			domain.ProjectMembershipCmd{
				UserID:    guest.ID.String(),
				ProjectID: project.ID.String(),
			},
		)

		accepted, err := s.AcceptInvitation(
			domain.NewContextWithUser(context.Background(), &guest),
			domain.ProjectMembershipCmd{
				ProjectID: projectMembership.ProjectID.String(),
			},
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if accepted.Role != "member" {
			t.Errorf("expected to be a member, want role %q, got %q", "member", accepted.Role)
		}
	})
}

func TestAcceptRequest(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewProjectMembershipService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		owner := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME17",
			FirstName: "FIRSTNAME17",
			LastName:  "LASTNAME17",
			Email:     "EMAIL17@email.com",
			Password:  "memaybeo",
		})
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME18",
			FirstName: "FIRSTNAME18",
			LastName:  "LASTNAME18",
			Email:     "EMAIL18@email.com",
			Password:  "huhuhuhu",
		})

		project := MustCreateProject(t, domain.NewContextWithUser(context.Background(), &owner), db, "TESTPROJECT2")

		projectMembership := MustRequest(
			t,
			domain.NewContextWithUser(context.Background(), &user),
			db,
			domain.ProjectMembershipCmd{
				UserID:    owner.ID.String(),
				ProjectID: project.ID.String(),
			},
		)

		accepted, err := s.AcceptRequest(
			domain.NewContextWithUser(context.Background(), &owner),
			domain.ProjectMembershipCmd{
				UserID:    user.ID.String(),
				ProjectID: projectMembership.ProjectID.String(),
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if accepted.Role != "member" {
			t.Errorf("expected to be a member, want role %q, got %q", "member", accepted.Role)
		}
	})
}

func TestDeleteMembership(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewProjectMembershipService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		owner := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME20",
			FirstName: "FIRSTNAME20",
			LastName:  "LASTNAME20",
			Email:     "EMAIL20@email.com",
			Password:  "memaybeo",
		})
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME21",
			FirstName: "FIRSTNAME21",
			LastName:  "LASTNAME21",
			Email:     "EMAIL21@email.com",
			Password:  "huhuhuhu",
		})

		ctxWithOwner := domain.NewContextWithUser(context.Background(), &owner)
		ctxWithUser := domain.NewContextWithUser(context.Background(), &user)

		project := MustCreateProject(t, ctxWithOwner, db, "PROJECT20")

		MustInvite(t, ctxWithOwner, db, domain.ProjectMembershipCmd{
			UserID:    user.ID.String(),
			ProjectID: project.ID.String(),
		})

		_, err := s.Delete(ctxWithOwner, domain.ProjectMembershipCmd{
			UserID:    user.ID.String(),
			ProjectID: project.ID.String(),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, _, err := s.Find(ctxWithUser, domain.ProjectMembershipFilter{
			ProjectID: nil,
			Offset:    0,
			Limit:     10,
		}); err == nil {
			t.Errorf("expected not found error")
		} else if domain.ErrorCode(err) != domain.ENOTFOUND || domain.ErrorMessage(err) != `no memberships found` {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestFindMemberships(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewProjectMembershipService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		owner := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME22",
			FirstName: "FIRSTNAME22",
			LastName:  "LASTNAME22",
			Email:     "EMAIL22@email.com",
			Password:  "memaybeo",
		})

		ctxWithOwner := domain.NewContextWithUser(context.Background(), &owner)

		mem1 := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME23",
			FirstName: "FIRSTNAME23",
			LastName:  "LASTNAME23",
			Email:     "EMAIL23@email.com",
			Password:  "huhuhuhu",
		})

		mem2 := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME24",
			FirstName: "FIRSTNAME24",
			LastName:  "LASTNAME24",
			Email:     "EMAIL24@email.com",
			Password:  "huhuhuhu",
		})

		project := MustCreateProject(t, ctxWithOwner, db, "TESTPROJECT22")

		MustInvite(t, ctxWithOwner, db, domain.ProjectMembershipCmd{
			UserID:    mem1.ID.String(),
			ProjectID: project.ID.String(),
		})

		MustInvite(t, ctxWithOwner, db, domain.ProjectMembershipCmd{
			UserID:    mem2.ID.String(),
			ProjectID: project.ID.String(),
		})

		projectID := project.ID.String()
		memberships, n, err := s.Find(ctxWithOwner, domain.ProjectMembershipFilter{
			ProjectID: &projectID,
			Offset:    0,
			Limit:     10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if n != 3 {
			t.Errorf("expect 3 membership in total, got %d", n)
		}

		if len(memberships) != 3 {
			t.Errorf("expect slice of 3 membership, got %d", len(memberships))
		}

		// memberships checks go here, i'm to burned out to do this shit now
	})
}

func MustInvite(
	tb testing.TB,
	ctx context.Context,
	db *postgres.DB,
	cmd domain.ProjectMembershipCmd,
) domain.ProjectMembership {
	projectMembership, err := postgres.NewProjectMembershipService(db, logrus.New()).Invite(ctx, cmd)
	if err != nil {
		tb.Fatalf("unexpected error: %v", err)
	}

	return projectMembership
}

func MustRequest(
	tb testing.TB,
	ctx context.Context,
	db *postgres.DB,
	cmd domain.ProjectMembershipCmd,
) domain.ProjectMembership {
	projectMembership, err := postgres.NewProjectMembershipService(db, logrus.New()).Request(ctx, cmd)
	if err != nil {
		tb.Fatalf("unexpected error: %v", err)
	}

	return projectMembership
}
