package user

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/app/validate"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/types"
)

type StubUserService struct {
	Users         map[string]types.User
	recordedID    uuid.UUID
	recordedEmail string
}

func (s *StubUserService) CreateUser(ctx context.Context, arg types.CreateUserCmd) (types.User, error) {

	dummyTime, _ := time.Parse(time.RFC3339, "2024-12-10T08:53:55+00:00")
	aUser := types.User{
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

func (s *StubUserService) ByID(ctx context.Context, id uuid.UUID) (types.User, error) {
	s.recordedID = id
	return types.User{}, nil
}

func (s *StubUserService) ByEmail(ctx context.Context, email string) (types.User, error) {

	s.recordedEmail = email
	return types.User{
		Username:  "validUser",
		Email:     email,
		FirstName: "Valid",
		LastName:  "",
	}, nil
}

func TestGetByEmail(t *testing.T) {
	email := "valid@email.com"
	service := &StubUserService{}

	type wantViewData struct {
		FirstName string
		Email     string
	}

	assertRendered := func(t testing.TB, got *goquery.Document, want wantViewData) {
		t.Helper()

		if got.Find(`form`).Length() == 0 {
			t.Error("expected form to be rendered, not found")
		}

		if !strings.Contains(got.Find(`h2`).Text(), want.FirstName) {
			t.Errorf("expected firstname to be rendered, got %s, want contains %s", got.Find(`h2`).Text(), want.FirstName)
		}

		if val, ok := got.Find(`form input[id="email"]`).Attr("value"); !ok {
			t.Error("expected email input to be rendered and has attribute value='<email>', not found")
		} else if val != want.Email {
			t.Errorf("expected email input has value=<email>, got %s want %s", val, want.Email)
		}

	}

	t.Run("it passes a valid email to service", func(t *testing.T) {
		e := echo.New()
		e.Validator = &validate.CustomValidator{Validator: validator.New()}

		f := createLoginCheckEmailFormParams(email)
		request := httptest.NewRequest(http.MethodPost, "/v1/users/login", strings.NewReader(f.Encode()))

		request.Header.Set("Content-Type", echo.MIMEApplicationForm)
		response := httptest.NewRecorder()
		h := NewHandler(service)

		c := e.NewContext(request, response)
		err := h.HandleLoginCheckEmail(c)
		if err != nil {
			t.Fatal(err)
		}

		if service.recordedEmail != email {
			t.Errorf("expected to pass a correct email, want %s got %s", email, service.recordedEmail)
		}

		doc, err := goquery.NewDocumentFromReader(response.Result().Body)
		if err != nil {
			t.Fatalf("failed to read template")
		}

		fmt.Println(doc.Text())

		want := wantViewData{
			Email:     "valid@email.com",
			FirstName: "Valid",
		}
		assertRendered(t, doc, want)
	})
}

func TestGetByID(t *testing.T) {

	uuidString := "7f173ec4-402d-4cd3-8446-0423771f972f"
	service := &StubUserService{
		Users: map[string]types.User{},
	}

	t.Run("it passes a valid id to service", func(t *testing.T) {

		e := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		response := httptest.NewRecorder()

		c := e.NewContext(request, response)
		c.SetParamNames("id")
		c.SetParamValues(uuidString)

		h := UserHandler{
			Service: service,
		}
		err := h.HandleUserGetByID(c)
		if err != nil {
			t.Fatal(err)
		}

		if service.recordedID != uuid.MustParse(uuidString) {
			t.Errorf("want %v got %v", uuid.MustParse(uuidString), service.recordedID)
		}
	})

	// t.Run("it returns error if user not found", func(t *testing.T) {
	// })
}

func TestCreateUser(t *testing.T) {

	service := &StubUserService{
		Users: map[string]types.User{},
	}

	type wantUserInfo struct {
		FirstName string
		LastName  string
	}

	assertResponse := func(t *testing.T, got *goquery.Document, want wantUserInfo) {
		t.Helper()
		wantFullName := fmt.Sprintf(`%s %s`, want.FirstName, want.LastName)

		if got.Find(`h2`).Length() == 0 {
			t.Errorf("expected a header, not found")
		}

		if got.Find(`p`).Length() == 0 {
			t.Errorf("expected welcome message, not found")
		}

		if link, ok := got.Find(`a`).Attr(`href`); !ok {
			t.Errorf("expected a link to login page, not found")
		} else if link != "/v1/login" {
			t.Errorf("expected link to '/v1/login', got '%s'", link)
		}

		if !strings.Contains(got.Find(`p`).Text(), wantFullName) {
			t.Errorf("expected user full name rendered, got %s", got.Find(`p`).Text())
		}
	}

	t.Run("it signup new user and render a response", func(t *testing.T) {
		// setup request
		f := createUserFormParams("TestUser", "Bob", "Ross", "bobr@email.com", "secretpassword")
		request, _ := http.NewRequest(http.MethodPost, "/v1/users/", strings.NewReader(f.Encode()))
		request.Header.Set("Content-Type", echo.MIMEApplicationForm)
		response := httptest.NewRecorder()

		// setup handler and inject fake
		e := echo.New()
		e.Validator = &validate.CustomValidator{Validator: validator.New()}
		h := UserHandler{
			Service: service,
		}

		// run handler
		c := e.NewContext(request, response)
		err := h.HandleUserCreate(c)
		if err != nil {
			t.Fatal(err)
		}

		if len(service.Users) != 1 {
			t.Errorf("expected %d users stored, got %d users stored", len(service.Users), 1)
		}

		doc, err := goquery.NewDocumentFromReader(response.Result().Body)
		if err != nil {
			t.Fatalf("failed to read template: %v", err)
		}

		want := wantUserInfo{
			FirstName: "Bob",
			LastName:  "Ross",
		}

		assertResponse(t, doc, want)
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
				e.Validator = &validate.CustomValidator{Validator: validator.New()}

				c := e.NewContext(request, response)
				h := UserHandler{
					Service: service,
				}
				h.HandleUserCreate(c)
				if response.Code != 400 {
					t.Errorf("Expected error code 400, got %d", response.Code)
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

func createLoginCheckEmailFormParams(email string) url.Values {
	f := make(url.Values)
	f.Set("email", email)
	return f
}
