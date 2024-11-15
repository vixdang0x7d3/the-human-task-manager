package app

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
		return template.Render(c, http.StatusBadRequest, components.AlertError("Invalid user info"))
	}

	user, err := h.Service.CreateUser(c.Request().Context(), types.CreateUserCmd(arg))
	if err != nil {
		return err
	}

	viewData := types.UserViewModel{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return template.Render(c, http.StatusAccepted, components.UserInfoPostSignup(viewData))
}

// DEPRECATED: I dont even know why i wrote this shit
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

func (h *UserHandler) HandleLoginCheckEmail(c echo.Context) error {

	type formData struct {
		Email string `form:"email" validate:"required,email"`
	}

	arg := formData{}
	if err := c.Bind(&arg); err != nil {
		return template.Render(c, http.StatusBadRequest, components.AlertError("bad request"))
	}

	if err := c.Validate(arg); err != nil {
		return template.Render(c, http.StatusBadRequest, components.AlertError("invalid email"))
	}

	user, err := h.Service.ByEmail(c.Request().Context(), arg.Email)
	if err != nil { // TODO: a more robust error handling
		return template.Render(c, http.StatusInternalServerError, components.AlertError("email not found or db error i can't tell :P"))
	}

	viewData := types.UserViewModel{
		FirstName: user.FirstName,
		Email:     user.Email,
	}

	return template.Render(c, http.StatusAccepted, components.UserLoginWithPassword(viewData))
}

func (h *UserHandler) HandleLoginCheckPassword(c echo.Context) error {
	type formData struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	arg := formData{}
	if err := c.Bind(&arg); err != nil {
		return err
	}

	if err := c.Validate(arg); err != nil {
		return template.Render(c, http.StatusBadRequest, components.AlertError("invalid password"))
	}

	user, err := h.Service.CheckPassword(c.Request().Context(), arg.Email, arg.Password)
	if err != nil {
		return template.Render(c, http.StatusBadRequest, components.AlertError("password is not correct"))
	}

	err = h.SessionManager.RenewToken(c.Request().Context())
	if err != nil {
		return template.Render(c, http.StatusInternalServerError, components.AlertError("fail to renew session"))
	}

	h.SessionManager.Put(c.Request().Context(), "userID", user.ID.String())
	h.SessionManager.Put(c.Request().Context(), "userEmail", user.Email)
	h.SessionManager.Put(c.Request().Context(), "userFistName", user.FirstName)

	return template.Render(c, http.StatusOK, components.PostLoginRedirectWithDelay())
}

func (h *UserHandler) HandleLoginRedirect(c echo.Context) error {
	c.Response().Header().Set("HX-Redirect", "/profile")
	return c.NoContent(http.StatusFound)
}

func (h *UserHandler) HandleLogout(c echo.Context) error {
	err := h.SessionManager.Destroy(c.Request().Context())
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "fail to destroy session")
	}

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusFound)
}