package user

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/template"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/template/components"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/types"
)

func (h *UserHandler) HandleUserCreate(c echo.Context) error {

	type formData struct {
		Username  string `form:"username" validate:"required"`
		FirstName string `form:"first_name" validate:"required"`
		LastName  string `form:"last_name" validate:"required"`
		Email     string `form:"email" validate:"required,email"`
		Password  string `form:"password" validate:"required"`
	}

	arg := formData{}
	if err := c.Bind(&arg); err != nil {
		return err
	}
	if err := c.Validate(arg); err != nil {
		return template.Render(c, http.StatusBadRequest, components.ErrorMessage("validation errors!"))
	}

	user, err := h.Service.CreateUser(c.Request().Context(), types.CreateUserCmd(arg))
	if err != nil {
		return err
	}

	data := types.ProfileViewModel{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return template.Render(c, http.StatusOK, components.UserInfoPostLogin(data))
}

func (h *UserHandler) HandleUserGetByID(c echo.Context) error {

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

	idString := c.Param("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Malformed id")
	}

	user, err := h.Service.ByID(c.Request().Context(), id)

	if err != nil {
		return err
	}

	if c.Request().Header.Get("Accept") == "text/html" {
		return echo.NewHTTPError(http.StatusBadRequest, "Unimplemented feature") // TODO: implement html rendering
	}

	return c.JSON(http.StatusAccepted, response(user))
}
