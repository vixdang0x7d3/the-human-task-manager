package app

import (
	"context"
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

type StubUserService struct {
	Users      map[string]domain.User
	recordedID uuid.UUID
}

func (s *StubUserService) CreateUser(ctx context.Context, arg domain.CreateUserParams) (domain.User, error) {

	dummyTime, _ := time.Parse(time.RFC3339, "2024-12-10T08:53:55+00:00")
	aUser := domain.User{
		ID:        uuid.MustParse("001af946-4f04-4dbf-a265-3be702667aea"),
		Username:  arg.Username,
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		Email:     arg.Email,
		SignupAt:  dummyTime,
		LastLogin: dummyTime,
	}
	s.Users[aUser.ID.String()] = aUser

	return aUser, nil
}

func (s *StubUserService) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	s.recordedID = id
	return domain.User{}, nil
}

func TestGetUser(t *testing.T) {

	uuidString := "7f173ec4-402d-4cd3-8446-0423771f972f"
	service := &StubUserService{
		Users: map[string]domain.User{},
	}

	t.Run("it passes a valid id to service", func(t *testing.T) {

		e := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/v1/users/", nil)
		response := httptest.NewRecorder()

		c := e.NewContext(request, response)
		c.SetParamNames("id")
		c.SetParamValues(uuidString)

		handler := UserHandler{
			Service: service,
		}
		err := handler.HandleGetUser(c)
		if err != nil {
			t.Fatal(err)
		}

		if service.recordedID != uuid.MustParse(uuidString) {
			t.Errorf("want %v got %v", uuid.MustParse(uuidString), service.recordedID)
		}
	})

	t.Run("it returns error if user not found", func(t *testing.T) {
	})
}

func TestCreateUser(t *testing.T) {

	service := &StubUserService{
		Users: map[string]domain.User{},
	}

	type WantAppUser struct {
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	assertAppUser := func(t *testing.T, got, want WantAppUser) {
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
			Service: service,
		}
		err := h.HandleCreateUser(c)
		if err != nil {
			t.Fatal(err)
		}

		if len(service.Users) != 1 {
			t.Errorf("expected %d users stored, got %d users stored", len(service.Users), 1)
		}

		var got WantAppUser
		if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}

		want := WantAppUser{
			Username:  "TestUser",
			FirstName: "Bob",
			LastName:  "Ross",
			Email:     "bobr@email.com",
		}

		assertAppUser(t, got, want)
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
				ErrorMsg: "validation error: invalid something",
			},

			{
				Description: "with missing password",
				Input: createUserFormParams(
					"TestUsername",
					"TestFirstName",
					"TestLastName",
					"abcxyz",
					"",
				),
				ErrorMsg: "validation error: invalid something",
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
					Service: service,
				}
				err := h.HandleCreateUser(c)
				if err == nil {
					t.Errorf("want error, got %s", err)
				}
			})
		}
	})
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
