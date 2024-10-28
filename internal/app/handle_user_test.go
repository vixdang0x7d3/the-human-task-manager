package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
)

const (
	DUMMY_TIME = "2024-12-10T08:53:55+00:00"
)

type StubUserService struct {
	Users map[string]domain.User
}

func (s *StubUserService) CreateUser(username, firstName, lastName, email, password string) (domain.User, error) {

	dummyTime, _ := time.Parse(time.RFC3339, DUMMY_TIME)

	aUser := domain.User{
		ID:        uuid.Nil,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		SignupAt:  dummyTime,
		LastLogin: dummyTime,
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

		// setup request
		f := createUserFormParams("TestUser", "Bob", "Ross", "bobr@email.com", "secretpassword")

		request, _ := http.NewRequest(http.MethodPost, "/v1/users/", strings.NewReader(f.Encode()))
		request.Header.Set("Content-Type", echo.MIMEApplicationForm)
		response := httptest.NewRecorder()

		e := echo.New()
		e.Validator = &internal.CustomValidator{Validator: validator.New()}

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

		want := AppUser{
			Username:  "TestUser",
			FirstName: "Bob",
			LastName:  "Ross",
			Email:     "bobr@email.com",
		}

		assertEqualUser(t, got, want)
	})

	t.Run("field validations", func(t *testing.T) {

		type testCase struct {
			Description string
			Input       url.Values
			ErrorMsg    string
		}

		for _, tc := range []testCase{
			{
				Description: "with invalid email",
				Input: createUserFormParams(
					"TestUsername",
					"TestFirstName",
					"TestLastName",
					"abcxyz",
					"TestPassword",
				),
				ErrorMsg: "validation error: invalid email",
			},

			{
				Description: "with missing username",
				Input: createUserFormParams(
					"",
					"TestFirstName",
					"TestLastName",
					"abcxyz",
					"TestPassword",
				),
				ErrorMsg: "validation error: invalid email",
			},

			{
				Description: "with missing password",
				Input: createUserFormParams(
					"",
					"TestFirstName",
					"TestLastName",
					"abcxyz",
					"TestPassword",
				),
				ErrorMsg: "validation error: invalid email",
			},
		} {
			t.Run(tc.Description, func(t *testing.T) {

				request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(tc.Input.Encode()))
				request.Header.Set("Content-Type", echo.MIMEApplicationForm)
				response := httptest.NewRecorder()

				e := echo.New()
				e.Validator = &internal.CustomValidator{Validator: validator.New()}

				c := e.NewContext(request, response)
				h := UserHandler{
					service: service,
				}
				err := h.HandleCreateUser(c)
				if err == nil {
					t.Errorf("want error, got %s", err)
				}
			})
		}
	})
}

func assertEqualUser(t *testing.T, got AppUser, want AppUser) {
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

func createUserFormParams(username, firstname, lastname, email, password string) url.Values {
	f := make(url.Values)
	f.Set("username", username)
	f.Set("first_name", firstname)
	f.Set("last_name", lastname)
	f.Set("email", email)
	f.Set("password", password)

	return f
}
