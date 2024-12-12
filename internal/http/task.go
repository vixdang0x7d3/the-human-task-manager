package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

func (s *Server) registerTaskRoutes(r *echo.Group) {
	r.GET("/tasks", s.handleTaskIndex)
	r.GET("/tasks/new", s.handleTaskNewShow)

	r.POST("/tasks/new", s.handleTaskNew)
}

func (s *Server) handleTaskIndex(c echo.Context) error {
	return nil
}

// handleTaskItem return task important attributes
// and aggregated data. Data returned by handleTaskItem
// is shown in Index screen
func (s *Server) handleTaskItem(c echo.Context) error {
	return nil
}

// handleTaskDetailShow returns a task object.
// Data returned by handleTaskDetailShow is
// shown in Task Update screen
func (s *Server) handleTaskDetailShow(c echo.Context) error {

	return render(c, http.StatusOK, pages.TaskList("/logout"))
}

func (s *Server) handleTaskDetailUpdate(c echo.Context) error {
	return nil
}

func (s *Server) handleTaskNewShow(c echo.Context) error {
	return nil
}

func (s *Server) handleTaskNew(c echo.Context) error {
	return nil
}
