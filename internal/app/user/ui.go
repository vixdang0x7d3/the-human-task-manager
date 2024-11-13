package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/template"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/template/pages"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/types"
)

func (h *UserHandler) HandleShowProfile(c echo.Context) error {

	user := types.UserViewModel{
		Username:  "bobr123",
		Email:     "bobr@email.com",
		FirstName: "Bob",
		LastName:  "Ross",
	}

	return template.Render(c, http.StatusOK, pages.Profile(user))
}

func (h *UserHandler) HandleShowSignup(c echo.Context) error {

	return template.Render(c, http.StatusOK, pages.Signup())
}

func (h *UserHandler) HandleShowLogin(c echo.Context) error {

	return template.Render(c, http.StatusOK, pages.Login())
}

func (h *UserHandler) HandleShowHome(c echo.Context) error {

	user := types.UserViewModel{
		Username:  "bobr123",
		Email:     "bobr@email.com",
		FirstName: "Bob",
		LastName:  "Ross",
	}

	return template.Render(c, http.StatusOK, pages.Home(user))
}
