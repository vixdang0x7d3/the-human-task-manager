package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/models"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/components"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

// registerAuthRoutes handle unauthenticated routes such as signup, login and oauth2.
// unauthenticated routes has prefix '/u'
func (s *Server) registerAuthRoutes(r *echo.Group) {
	r.GET("/login", s.handleLoginShow)
	r.GET("/signup", s.handleSignupShow)

	r.POST("/signup", s.handleSignup)
	r.POST("/login/email", s.handleLoginEmail)
	r.POST("/login/password", s.handleLoginPassword)
	r.GET("/login/success", s.handleLoginSuccess)
}

func (s *Server) handleLoginShow(c echo.Context) error {
	return render(c, http.StatusOK, pages.LoginEmail("/u/login/email", "/u/signup"))
}

func (s *Server) handleSignupShow(c echo.Context) error {
	return render(c, http.StatusOK, pages.Signup("/u/signup"))
}

func (s *Server) handleSignup(c echo.Context) error {
	type formValues struct {
		Username  string `form:"username" validate:"required"`
		FirstName string `form:"first_name" validate:"required"`
		LastName  string `form:"last_name" validate:"required"`
		Email     string `form:"email" validate:"required,email"`
		Password  string `form:"password" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusBadRequest, components.AlertError("internal error"))
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error("form validation error", err)
		return render(c, http.StatusBadRequest, components.AlertError("invalid signup info"))
	}

	user, err := s.UserService.Create(c.Request().Context(), domain.CreateUserCmd(form))
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
	}

	m := models.UserView{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return render(c, http.StatusOK, components.UserInfoPostSignup(m))
}

func (s *Server) handleLoginEmail(c echo.Context) error {
	type formValues struct {
		Email string `form:"email" validate:"required,email"`
	}

	form := formValues{}

	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusBadRequest, components.AlertError("invalid form data"))
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error("form validation error", err)
		return c.HTML(http.StatusBadRequest, "invalid email")
	}

	user, err := s.UserService.ByEmail(c.Request().Context(), form.Email)
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
	}

	m := models.UserView{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return render(c, http.StatusOK, components.LoginPassword(m, "/u/login/password"))
}

func (s *Server) handleLoginPassword(c echo.Context) error {

	type formValues struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid form data")
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid password")
	}

	user, err := s.UserService.ByEmailWithPassword(c.Request().Context(), form.Email, form.Password)
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
	}

	if err := s.sessions.RenewToken(c.Request().Context()); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusInternalServerError, "internal error")
	}
	s.sessions.Put(c.Request().Context(), "userID", user.ID.String())

	return render(c, http.StatusOK, components.LoginAlertRedirect("login successful!", "/u/login/success"))
}

func (s *Server) handleLoginSuccess(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
