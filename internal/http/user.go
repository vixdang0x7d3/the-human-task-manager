package http

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/models"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/components"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

func (s *Server) registerUserRoutes(r *echo.Group) {
	r.GET("/index", s.handleIndexShow)
	r.GET("/profile", s.handleProfileShow)

	r.DELETE("/logout", s.handleLogout)
}

func (s *Server) handleIndexShow(c echo.Context) error {

	user := domain.UserFromContext(c.Request().Context())
	if user == nil {
		c.Logger().Error("no user in current session")
		return echo.NewHTTPError(http.StatusBadRequest, "no user in login session")
	}

	m := models.UserView{
		FirstName: user.FirstName,
	}

	return render(c, http.StatusOK, pages.Index(m, "/logout"))
}

func (s *Server) handleProfileShow(c echo.Context) error {

	user := domain.UserFromContext(c.Request().Context())

	m := models.UserView{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return render(c, http.StatusOK, pages.Profile(m, "/change-info", "/logout"))
}

func (s *Server) handleChangeProfileInfo(c echo.Context) error {

	type formValues struct {
		Username  string `form:"username" validate:"required"`
		FirstName string `form:"first_name" validate:"required"`
		LastName  string `form:"last_name" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid form data")
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid value")
	}

	m := models.UserView{
		Username:  form.Username,
		FirstName: form.FirstName,
		LastName:  form.LastName,
	}

	return render(c, http.StatusOK, components.ChangeInfoForm(m, "/change-info-save"))
}

// TODO: flash message support
func (s *Server) handleLogout(c echo.Context) error {
	err := s.sessions.Destroy(c.Request().Context())
	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, "internal error")
	}

	c.Response().Header().Set("HX-Redirect", "/u/login")
	return c.NoContent(http.StatusOK)
}
