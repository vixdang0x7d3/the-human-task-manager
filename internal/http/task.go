package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/generic"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/models"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/components"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

func (s *Server) registerTaskRoutes(r *echo.Group) {
	r.GET("/tasks", s.handleTaskIndex)
	r.GET("/tasks/new", s.handleTaskNewShow)

	r.POST("/tasks/new", s.handleTaskNew)
}

func (s *Server) handleTaskIndex(c echo.Context) error {

	state := "started"
	taskItems, total, err := s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		State:  &state,
		Offset: 0,
		Limit:  100,
	})

	// handle pagination
	_ = total

	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
		case domain.EINTERNAL:
			c.Logger().Error(domain.ErrorMessage(err))
			return c.HTML(http.StatusInternalServerError, "internal error")
		}
	}

	return render(c, http.StatusOK, pages.TaskIndex(
		generic.Map(taskItems, toTaskItemView),
		"/logout",
	))
}

func toTaskItemView(taskItem domain.TaskItem) models.TaskItemView {

	const timeFormat = "2006/01/02 15:04:05"

	var (
		deadline string = "none"
		schedule string = "none"
		wait     string = "none"
		end      string = "none"
	)

	if !taskItem.Deadline.IsZero() {
		deadline = taskItem.Deadline.Format(timeFormat)
	}

	if !taskItem.Schedule.IsZero() {
		schedule = taskItem.Schedule.Format(timeFormat)
	}

	if !taskItem.Wait.IsZero() {
		deadline = taskItem.Wait.Format(timeFormat)
	}

	if !taskItem.End.IsZero() {
		deadline = taskItem.End.Format(timeFormat)
	}

	return models.TaskItemView{
		ID:             taskItem.ID.String(),
		Description:    taskItem.Description,
		UserID:         taskItem.UserID.String(),
		Username:       taskItem.Username,
		CompleteBy:     taskItem.CompletedBy.String(),
		CompleteByName: taskItem.CompletedByName,
		ProjectID:      taskItem.ProjectID.String(),
		ProjectTitle:   taskItem.ProjectTitle,

		Priority: taskItem.Priority,
		State:    taskItem.State,

		Deadline: deadline,
		Schedule: schedule,
		Wait:     wait,
		Create:   taskItem.Create.Format(timeFormat),
		End:      end,
		Urgency:  strconv.FormatFloat(taskItem.Urgency, 'f', 2, 64),

		Tags: taskItem.Tags,
	}
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
func (s *Server) handleTaskUpdateShow(c echo.Context) error {
	return nil
}

func (s *Server) handleTaskUpdate(c echo.Context) error {
	return nil
}

func (s *Server) handleTaskNewShow(c echo.Context) error {

	priorities := []string{
		domain.TaskPriorityH,
		domain.TaskPriorityM,
		domain.TaskPriorityL,
	}

	projects, total, err := s.ProjectService.Find(c.Request().Context(), domain.ProjectFilter{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return c.HTML(http.StatusInternalServerError, "internal error")
		}
	}

	if total != len(projects) {
		c.Logger().Warn("truncated slice of projects ", len(projects), total)
	}

	return render(c, http.StatusOK, pages.TaskNew(
		[]string{"next", "school", "personal"},
		priorities,
		generic.Map(projects, toProjectView),
		"/logout",
	))
}

func toProjectView(project domain.Project) models.ProjectView {
	return models.ProjectView{
		Title:  project.Title,
		ID:     project.ID.String(),
		UserID: project.UserID.String(),
	}
}

func (s *Server) handleTaskNew(c echo.Context) error {

	type formValues struct {
		Description string `form:"description" validate:"required"`
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
		return render(c, http.StatusBadRequest, components.AlertError("invalid form data"))
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusBadRequest, components.AlertError("invalid task info"))
	}

	tags, err := parseTagsJSON([]byte(form.TagsJSON))
	if err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusBadRequest, components.AlertError("invalid tags format"))
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
			c.Logger().Error(domain.ErrorMessage(err))
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(domain.ErrorMessage(err))
			return render(c, http.StatusBadRequest, components.AlertError("internal error"))
		}
	}

	c.Logger().Info("create task success ", task.ID)

	c.Response().Header().Set("HX-Redirect", "/tasks")
	return c.NoContent(http.StatusFound)
}
