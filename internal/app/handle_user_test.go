package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
)

const (
	MOCK_TIMESTAMP = "2024-12-10T08:53:55+00:00"
)

type StubUserService struct {
	Users map[string]domain.User
}

func (s *StubUserService) CreateUser(username, firstName, lastName, email, password string) (domain.User, error) {

	aUser := domain.User{
		ID:        uuid.New(),
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		SignupAt:  time.Now(),
		LastLogin: time.Now(),
	}
	s.Users[aUser.ID.String()] = aUser

	return aUser, nil
}

func (s *StubUserService) GetUser(userID string) (domain.User, error) {
	return domain.User{}, nil
}

func TestCreateUser(t *testing.T) {

	service := &StubUserService{
		map[string]domain.User{},
	}

	t.Run("it records a new user", func(t *testing.T) {

		want := CreateUserParam{
			Username:  "TestUser",
			FirstName: "Bob",
			LastName:  "Ross",
			Email:     "bobr@email.company",
			Password:  "secretpassword",
		}

		e := echo.New()
		request, _ := http.NewRequest(http.MethodPost, "/v1/users/", nil)
		response := httptest.NewRecorder()

		c := e.NewContext(request, response)
		h := UserHandler{
			service: service,
		}
		h.HandleCreateUser(c)

		if len(service.Users) != 1 {
			t.Errorf("expected %d users stored, got %d users stored", len(service.Users), 1)
		}
	})
}
