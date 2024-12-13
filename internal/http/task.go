package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/models"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

func (s *Server) registerTaskRoutes(r *echo.Group) {
	r.GET("/tasks", s.handleTaskIndex)
	r.GET("/tasks/new", s.handleTaskNewShow)

	r.POST("/tasks/new", s.handleTaskNew)
}

func (s *Server) handleTaskIndex(c echo.Context) error {

	return render(c, http.StatusOK, pages.TaskIndex("/logout"))
}

// handleTaskItem return task important attributes
// and aggregated data. Data returned by handleTaskItem
// is shown in Index screen (i might not need this end-point)
func (s *Server) handleTaskItem(c echo.Context) error {
	return nil
}

// handleTaskDetailShow returns a task object.
// Data returned by handleTaskDetailShow is
// shown in Task Update screen
func (s *Server) handleTaskDetailShow(c echo.Context) error {
	return nil
}

func (s *Server) handleTaskDetailUpdate(c echo.Context) error {
	return nil
}

func (s *Server) handleTaskNewShow(c echo.Context) error {

	priorities := []string{
		domain.TaskPriorityH,
		domain.TaskPriorityM,
		domain.TaskPriorityL,
	}

	projects := []models.ProjectView{
		{
			Title: "Project A",
			ID:    "project-a-id",
		},
		{
			Title: "Project B",
			ID:    "project-b-id",
		},
	}

	return render(c, http.StatusOK, pages.TaskNew(
		[]string{"your", "mom", ">:("},
		priorities,
		projects,
		"/logout",
	))
}

func (s *Server) handleTaskNew(c echo.Context) error {

	type formValues struct {
		Description string `form:"description"`
		Deadline    string `form:"deadline"`
		Schedule    string `form:"schedule"`
		Wait        string `form:"wait"`
		TagsJSON    string `form:"tags"`
		Priority    string `form:"priority"`
		ProjectID   string `form:"project_id"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid form data")
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid task info")
	}

	tags, err := parseTagsJSON([]byte(form.TagsJSON))
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid tags json string")
	}

	task, err := s.TaskService.Create(c.Request().Context(), domain.CreateTaskCmd{
		Description: form.Description,
		Deadline:    form.Deadline,
		Schedule:    form.Schedule,
		Wait:        form.Wait,
		Tags:        tags,
		Priority:    form.Priority,
		ProjectID:   form.ProjectID,
	})
	if err != nil {
		switch domain.ErrorMessage(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
		case domain.EINTERNAL:
			c.Logger().Error(domain.ErrorMessage(err))
			return c.HTML(http.StatusInternalServerError, "internal error")
		}
	}

	c.Logger().Info("create task success ", task.ID)

	c.Response().Header().Set("HX-Redirect", "/tasks")
	return c.NoContent(http.StatusFound)
}
