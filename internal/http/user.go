package http

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/models"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

func (s *Server) registerUserRoutes(r *echo.Group) {
	r.GET("/index", s.handleIndexShow)
	r.GET("/profile", s.handleProfileShow)

	r.DELETE("/logout", s.handleLogout)
}

func (s *Server) handleIndexShow(c echo.Context) error {

	sessionFirstName := s.sessions.GetString(c.Request().Context(), "userFirstName")
	if sessionFirstName == "" {
		return c.HTML(http.StatusUnauthorized, "unauthorized user")
	}

	m := models.UserView{
		FirstName: sessionFirstName,
	}

	return render(c, http.StatusOK, pages.Index(m, "/logout"))
}

func (s *Server) handleProfileShow(c echo.Context) error {

	id := s.sessions.GetString(c.Request().Context(), "userID")

	user, err := s.UserService.ByID(c.Request().Context(), id)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
	}

	m := models.UserView{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return render(c, http.StatusOK, pages.Profile(m, "/logout"))
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
