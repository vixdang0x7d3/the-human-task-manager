package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/models"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

func (s *Server) registerCalendarRoutes(r *echo.Group) {
	r.GET("/calendar", s.handleCalendarShow)
}

func (s *Server) handleCalendarShow(c echo.Context) error {

	m := []models.TaskView{
		{
			Title:    "task 1",
			Schedule: "2024-11-26",
		},
		{
			Title:    "task 2",
			Schedule: "2024-11-28",
		},
	}

	return render(c, http.StatusOK, pages.Calendar(m, "/logout"))
}
