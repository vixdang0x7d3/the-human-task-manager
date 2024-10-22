package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

	mockTime, _ := time.Parse(time.RFC3339, MOCK_TIMESTAMP)

	aUser := domain.User{
		ID:        uuid.Nil,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		SignupAt:  mockTime,
		LastLogin: mockTime,
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
		params := CreateUserParam{
			Username:  "TestUser",
			FirstName: "Bob",
			LastName:  "Ross",
			Email:     "bobr@email.company",
			Password:  "secretpassword",
		}

		mockTime, _ := time.Parse(time.RFC3339, MOCK_TIMESTAMP)
		want := AppUser{
			ID:        uuid.Nil,
			Username:  params.Username,
			FirstName: params.FirstName,
			LastName:  params.LastName,
			Email:     params.Email,
			SignupAt:  mockTime,
			LastLogin: mockTime,
		}

		e := echo.New()

		f := make(url.Values)
		f.Set("username", params.Username)
		f.Set("first_name", params.FirstName)
		f.Set("last_name", params.LastName)
		f.Set("email", params.Email)
		f.Set("password", params.Password)

		request, _ := http.NewRequest(http.MethodPost, "/v1/users/", strings.NewReader(f.Encode()))
		request.Header.Set("Content-Type", echo.MIMEApplicationForm)

		response := httptest.NewRecorder()

		c := e.NewContext(request, response)
		h := UserHandler{
			service: service,
		}
		err := h.HandleCreateUser(c)
		if err != nil {
			t.Fatal(err)
		}

		if len(service.Users) != 1 {
			t.Errorf("expected %d users stored, got %d users stored", len(service.Users), 1)
		}

		var got AppUser
		if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}

		// go time.Time can not be compared normally :D
		assertEqualUser(t, got, want)
	})
}

func assertEqualUser(t *testing.T, got, want AppUser) {
	t.Helper()
	if got.ID != want.ID {
		t.Errorf("got %s want %s", got.ID, want.ID)
	}
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
	if !got.LastLogin.Equal(want.LastLogin) {
		t.Errorf("got %s want %s", got.LastLogin, want.LastLogin)
	}
	if !got.SignupAt.Equal(want.SignupAt) {
		t.Errorf("got %s want %s", got.SignupAt, want.SignupAt)
	}
}
