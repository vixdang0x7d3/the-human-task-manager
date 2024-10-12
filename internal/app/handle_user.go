package app

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
)

type UserHandler struct {
	Service UserService
}

func (h *UserHandler) HandleCreateUser(c echo.Context) error {
	payload := &UserFormPayload{}

	if err := c.Bind(payload); err != nil {
		return fmt.Errorf("error when parsing form payload: %v", err)
	}

	user, err := h.Service.CreateUser(
		payload.Username,
		payload.FirstName,
		payload.LastName,
		payload.Email,
		payload.Password,
	)

	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return c.JSON(http.StatusOK, toAppUser(user))
}

type UserService interface {
	CreateUser(username, firstname, lastname, email, password string) (domain.User, error)
	GetUserByID(userID string) (domain.User, error)
}
