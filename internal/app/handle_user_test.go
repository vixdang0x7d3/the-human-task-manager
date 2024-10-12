package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
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

type userServiceStub struct{}

func (stub *userServiceStub) CreateUser(username, firstname, lastname, email, password string) (domain.User, error) {
	mockTimeValue, _ := time.Parse(time.RFC3339, MOCK_TIMESTAMP)
	return domain.User{
		ID:        uuid.Nil,
		Username:  username,
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		SignupAt:  mockTimeValue,
		LastLogin: mockTimeValue,
	}, nil
}

func (stub *userServiceStub) GetUserByID(userID string) (domain.User, error) {
	return domain.User{}, nil
}

func TestCreateUser(t *testing.T) {

	t.Run("return user's info ", func(t *testing.T) {
		formData := map[string]string{
			"username":   "imbobr",
			"first_name": "Bob",
			"last_name":  "Ross",
			"email":      "bob.ross@email.com",
			"password":   "bobsecret",
		}

		f := make(url.Values)
		f.Set("username", formData["username"])
		f.Set("first_name", formData["first_name"])
		f.Set("last_name", formData["last_name"])
		f.Set("email", formData["email"])
		f.Set("password", formData["password"])

		timeMockValue, _ := time.Parse(time.RFC3339, MOCK_TIMESTAMP)
		wantedAppUser := AppUser{
			ID:        uuid.Nil,
			Username:  formData["username"],
			FirstName: formData["first_name"],
			LastName:  formData["last_name"],
			Email:     formData["email"],
			SignupAt:  timeMockValue,
			LastLogin: timeMockValue,
		}
		userJSON, _ := json.Marshal(&wantedAppUser)

		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "/users/", strings.NewReader(f.Encode()))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		response := httptest.NewRecorder()

		c := e.NewContext(request, response)

		handler := &UserHandler{
			Service: &userServiceStub{},
		}
		err := handler.HandleCreateUser(c)
		if err != nil {
			t.Fatal(err)
		}

		if response.Code != http.StatusOK {
			t.Errorf("got status code %d want %d", response.Code, http.StatusOK)
		}

		if reflect.DeepEqual(response.Body.String(), string(userJSON)) {
			t.Errorf("\ngot %s want %s", response.Body.String(), userJSON)
		}
	})
}
