package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/types"
)

type StubUserStore struct {
	CreateCount int
	Users       map[uuid.UUID]database.User
}

func (s *StubUserStore) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	s.CreateCount++
	s.Users[arg.ID] = database.User(arg)
	return s.Users[arg.ID], nil
}

func (s *StubUserStore) ByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	return database.User{}, nil
}

func (s *StubUserStore) ByEmail(ctx context.Context, email string) (database.User, error) {
	for _, u := range s.Users {
		if u.Email == email {
			return u, nil
		}
	}
	return database.User{}, errors.New("No user")
}

func TestCreateUser(t *testing.T) {
	s := &StubUserStore{
		CreateCount: 0,
		Users:       map[uuid.UUID]database.User{},
	}

	t.Run("it creates a user", func(t *testing.T) {
		plainTextPassword := "secretpassword"

		arg := types.CreateUserCmd{
			Username:  "TestUSer",
			FirstName: "Bob",
			LastName:  "Ross",
			Email:     "bobr@email.company",
			Password:  plainTextPassword,
		}

		wantDomainUser := types.User{
			ID:        uuid.New(),
			Username:  arg.Username,
			FirstName: arg.FirstName,
			LastName:  arg.LastName,
			Email:     arg.Email,
			SignupAt:  time.Now(),
			LastLogin: time.Now(),
		}

		core := UserCore{s}
		gotDomainUser, err := core.CreateUser(context.Background(), arg)
		if err != nil {
			t.Fatal(err)
		}

		if len(s.Users) != 1 {
			t.Errorf("expected 1 user saved got %d", len(s.Users))
		}

		if _, ok := s.Users[gotDomainUser.ID]; !ok {
			t.Errorf("user is not properly stored to db")
		}

		if !checkPassword(plainTextPassword, s.Users[gotDomainUser.ID].Password) {
			t.Errorf("password is not properly hashed")
		}

		assertDomainUser(t, gotDomainUser, wantDomainUser)
	})
}

func TestCheckPassword(t *testing.T) {

	hashedPassword, _ := hashPassword("secret")
	ignoreMe, _ := hashPassword("ignoreme")
	s := &StubUserStore{
		Users: map[uuid.UUID]database.User{
			uuid.MustParse("ce25c916-f064-4e2f-a6db-bb15b87f599a"): {
				ID:       uuid.MustParse("ce25c916-f064-4e2f-a6db-bb15b87f599a"),
				Username: "TestUser",
				Email:    "test@email.com",
				Password: hashedPassword,
			},

			uuid.MustParse("522d0c6b-d922-4d2e-81d6-2e69fe7a8906"): {
				ID:       uuid.MustParse("522d0c6b-d922-4d2e-81d6-2e69fe7a8906"),
				Username: "AnotherUser",
				Email:    "another@email.com",
				Password: ignoreMe,
			},
		},
	}

	t.Run("happy path", func(t *testing.T) {

		email := "test@email.com"
		password := "secret"

		c := NewUserCore(s)
		user, err := c.CheckPassword(context.Background(), email, password)
		if err != nil {
			t.Errorf("should not have error, unexpected error: %v", err)
		}

		if user.Email != email {
			t.Errorf("return wrong user's email")
		}
	})
}

func assertDomainUser(t *testing.T, got, want types.User) {
	t.Helper()
	if got.Username != want.Username {
		t.Errorf("got %s want %s", got.Username, want.Username)
	}
	if got.FirstName != want.FirstName {
		t.Errorf("got %s want %s", got.FirstName, want.FirstName)
	}
	if got.LastName != want.LastName {
		t.Errorf("got %s want %s", got.LastName, want.LastName)
	}
	if got.Email != want.Email {
		t.Errorf("got %s want %s", got.Email, want.Email)
	}
}
