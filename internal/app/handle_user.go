package app

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
)

type UserHandler struct {
	service UserService
}

func (h *UserHandler) HandleCreateUser(c echo.Context) error {

	type request struct {
		Username  string `form:"username" validate:"required"`
		FirstName string `form:"first_name" validate:"required"`
		LastName  string `form:"last_name" validate:"required"`
		Email     string `form:"email" validate:"required,email"`
		Password  string `form:"password" validate:"required"`
	}

	// used for testing purpose
	type response struct {
		ID        uuid.UUID `json:"id"`
		Username  string    `json:"username"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		SignupAt  time.Time `json:"signup_at"`
		LastLogin time.Time `json:"last_login"`
	}

	arg := &request{}
	if err := c.Bind(arg); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(arg); err != nil {
		return err
	}

	user, err := h.service.CreateUser(
		arg.Username,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Password,
	)
	if err != nil {
		return err
	}

	if c.Request().Header.Get("Accept") == "text/html" {
		return echo.ErrInternalServerError // TODO: implement html rendering
	}

	return c.JSON(http.StatusAccepted, response(user))
}

type UserService interface {
	CreateUser(username, firstName, lastName, email, password string) (domain.User, error)
	GetUser(userID string) (domain.User, error)
}
