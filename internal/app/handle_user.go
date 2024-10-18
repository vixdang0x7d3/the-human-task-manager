package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
)

type UserHandler struct {
	service UserService
}

func (h *UserHandler) HandleCreateUser(c echo.Context) error {
	arg := &CreateUserParam{}
	if err := c.Bind(arg); err != nil {
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

	return c.JSON(http.StatusAccepted, AppUser(user))
}

type UserService interface {
	CreateUser(username, firstName, lastName, email, password string) (domain.User, error)
	GetUser(userID string) (domain.User, error)
}
