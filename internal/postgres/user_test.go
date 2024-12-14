package postgres_test

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres"
)

func TestCreateUser(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db := MustOpenDB(t, context.Background())
		defer CloseDB(t, db)

		s := postgres.NewUserService(db, logrus.New())

		cmd := domain.CreateUserCmd{
			Username:  "nhunghongnhung",
			FirstName: "Thi Hong Nhung",
			LastName:  "Nguyen",
			Email:     "hongnhungnguyen1304@gmail.com",
		}

		user, err := s.Create(context.Background(), cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if user.Username != cmd.Username {
			t.Errorf("mismatch %q != %q", user.Username, cmd.Username)
		}

		if user.FirstName != cmd.FirstName {
			t.Errorf("mismatch %q != %q", user.FirstName, cmd.FirstName)
		}

		if user.LastName != cmd.LastName {
			t.Errorf("mismatch %q != %q", user.LastName, cmd.LastName)
		}

		if user.Email != cmd.Email {
			t.Errorf("mismatch %q != %q", user.Email, cmd.Email)
		}

		if user.SignupAt.IsZero() {
			t.Errorf("expected signup at")
		}

		if user.LastLogin.IsZero() {
			t.Errorf("expected last login")
		}
	})

	t.Run("ErrEmailExists", func(t *testing.T) {
		db := MustOpenDB(t, context.Background())
		defer CloseDB(t, db)

		s := postgres.NewUserService(db, logrus.New())

		_, err := s.Create(context.Background(), domain.CreateUserCmd{
			Username:  "USERNAME",
			FirstName: "FIRSTNAME",
			LastName:  "LASTNAME",
			Email:     "existed@email.com",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		user, err := s.Create(context.Background(), domain.CreateUserCmd{
			Username:  "USERNAME",
			FirstName: "FIRSTNAME",
			LastName:  "LASTNAME",
			Email:     "existed@email.com",
		})

		_ = user
		if err == nil {
			t.Error("expected error, not found")
		} else if domain.ErrorCode(err) != domain.ECONFLICT ||
			domain.ErrorMessage(err) != `createUser: this email exists` {
			t.Errorf(`unexpected error: %#v`, err)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	s := postgres.NewUserService(db, logrus.New())

	t.Run("OK", func(t *testing.T) {
		user := MustCreateUser(t, context.Background(), db, domain.CreateUserCmd{
			Username:  "USERNAME27",
			FirstName: "FIRSTNAME27",
			LastName:  "LASTNAME27",
			Email:     "EMAIL27@email.com",
		})

		ctxWithUser := domain.NewContextWithUser(context.Background(), &user)

		username := "steveee"
		firstName := "Steve"
		lastName := "Dep Trai"
		email := "stevedt@email.com"
		password := "stevebibede"

		_, err := s.Update(ctxWithUser, domain.UpdateUserCmd{
			Username:  &username,
			FirstName: &firstName,
			LastName:  &lastName,
			Email:     &email,
			Password:  &password,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		updated, err := s.ByID(ctxWithUser, user.ID.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if updated.Username != username {
			t.Errorf("username mismatch: %q != %q", updated.Username, username)
		}

		if updated.FirstName != firstName {
			t.Errorf("first name mismatch: %q != %q", updated.FirstName, firstName)
		}

		if updated.LastName != lastName {
			t.Errorf("last name mismatch: %q != %q", updated.LastName, lastName)
		}

		if updated.Email != email {
			t.Errorf("email mismatch: %q != %q", updated.Email, email)
		}

		if _, err := s.ByEmailWithPassword(context.Background(), email, password); err != nil {
			t.Errorf("password mismatch, got error: %v", err)
		}
	})
}

func TestUserByEmail(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)
	t.Skip("no test")
}

func TestUserByEmailWithPassword(t *testing.T) {

	db := MustOpenDB(t, context.Background())
	defer CloseDB(t, db)

	arg := domain.CreateUserCmd{
		Username:  "USERNAME",
		FirstName: "FIRSTNAME",
		LastName:  "LASTNAME",

		Email:    "bob@email.com",
		Password: "secretpassword",
	}

	s := postgres.NewUserService(db, logrus.New())

	_, err := s.Create(context.Background(), arg)
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	t.Run("OK", func(t *testing.T) {

		user, err := s.ByEmailWithPassword(context.Background(), arg.Email, arg.Password)
		if err != nil {
			t.Fatalf("should have no error: %#v", err)
		}

		if user.Email != arg.Email {
			t.Errorf("email mismatch: %q != %q", user.Email, arg.Email)
		}
	})

	t.Run("ErrNotFound", func(t *testing.T) {

		_, err := s.ByEmailWithPassword(context.Background(), "invalid@email.com", "notimportant")
		if err == nil {
			t.Error("expected error, not found")
		} else if domain.ErrorCode(err) != domain.ENOTFOUND ||
			domain.ErrorMessage(err) != `userByEmailWithPassword: email not found` {
			t.Errorf("unexpected error: %#v", err)
		}

	})

	t.Run("ErrWrongPassword", func(t *testing.T) {
		_, err := s.ByEmailWithPassword(context.Background(), arg.Email, "wrongpassword")
		if err == nil {
			t.Error("expected error, not found")
		} else if domain.ErrorCode(err) != domain.EUNAUTHORIZED ||
			domain.ErrorMessage(err) != `userByEmailWithPassword: wrong password` {

			t.Errorf("unexpected error: %#v", err)
		}
	})
}

func MustCreateUser(tb testing.TB, ctx context.Context, db *postgres.DB, cmd domain.CreateUserCmd) domain.User {
	tb.Helper()
	user, err := postgres.NewUserService(db, logrus.New()).Create(ctx, cmd)
	if err != nil {
		tb.Fatal(err)
	}
	return user
}
