package http

import (
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/generic"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/components"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

func (s *Server) registerTaskRoutes(r *echo.Group) {
	r.GET("/tasks", s.handleTaskIndex)
	r.GET("/tasks/new", s.handleTaskNewShow)
	r.GET("/tasks/update/:id", s.handleTaskUpdateShow)

	r.POST("/tasks/find", s.handleTaskFind)
	r.POST("/tasks/new", s.handleTaskNew)
	r.POST("/tasks/update/:id", s.handleTaskUpdate)
	r.POST("/tasks/set-project/:id", s.handleTaskSetProject)
	r.POST("/tasks/complete/:id", s.handleTaskComplete)
	r.POST("/tasks/start/:id", s.handleTaskStart)
	r.DELETE("/tasks/delete/:id", s.handleTaskDelete)

}

const limit = 5 // limit set the max number of tasks shown in a page

func (s *Server) handleTaskIndex(c echo.Context) error {

	var pageOffset int
	if qParam := c.QueryParam("pageOffset"); qParam != "" {
		parsed, err := strconv.ParseInt(qParam, 10, 64)
		if err == nil {
			pageOffset = int(parsed)
		}

	}

	state := "started"
	taskItems, total, err := s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		State:  &state,
		Offset: pageOffset * limit,
		Limit:  limit,
	})

	pageTotal := int(math.Ceil(float64(total) / float64(limit)))
	taskTotal := total

	projects, total, err := s.ProjectService.Find(c.Request().Context(), domain.ProjectFilter{
		Offset: 0,
		Limit:  100,
	})

	if total > 100 {
		c.Logger().Warnf("truncated projects list, total: %d, show: 100", total)
	}

	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
		case domain.EINTERNAL:
			c.Logger().Error(domain.ErrorMessage(err))
			return c.HTML(http.StatusInternalServerError, "internal error")
		}
	}

	trigger := c.Request().Header.Get("HX-Trigger")
	if trigger == "prev-btn" || trigger == "next-btn" {
		return render(c, http.StatusOK, components.TaskList(
			generic.Map(taskItems, toTaskItemView),
			generic.Map(projects, toProjectView),
			pageOffset,
			pageTotal,
			taskTotal,
		))
	}

	return render(c, http.StatusOK, pages.TaskIndex(
		generic.Map(taskItems, toTaskItemView),
		generic.Map(projects, toProjectView),
		pageOffset,
		pageTotal,
		taskTotal,
		"/logout",
	))
}

func (s *Server) handleTaskFind(c echo.Context) error {

	type formValues struct {
		Query    string `form:"query"`
		State    string `form:"state"`
		Priority string `form:"priority"`
		Days     string `form:"days"`
		Months   string `form:"months"`
	}

	form := formValues{}

	if err := c.Bind(&form); err != nil {
		return c.HTML(http.StatusBadRequest, "invalid form data")
	}

	var (
		query    *string
		state    *string
		priority *string
		days     *int64
		months   *int64
	)

	if form.Query != "" {
		query = &form.Query
	}

	if form.State != "" {
		state = &form.State
	}

	if form.Priority != "" {
		priority = &form.Priority
	}

	if form.Days != "" {
		val, err := strconv.ParseInt(form.Days, 10, 64)
		if err == nil && val != 0 {
			days = &val
		}
		c.Logger().Debugf("parse month has problem: %v", err)
	}

	if form.Months != "" {
		val, err := strconv.ParseInt(form.Months, 10, 64)
		if err == nil && val != 0 {
			months = &val

		}
		c.Logger().Debugf("parse month has problem: %v", err)
	}

	var pageOffset int
	if qParam := c.QueryParam("pageOffset"); qParam != "" {
		parsed, err := strconv.ParseInt(qParam, 10, 64)
		if err == nil {
			pageOffset = int(parsed)
		}
	}

	taskItems, total, err := s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		Q:        query,
		State:    state,
		Priority: priority,
		Days:     days,
		Months:   months,
		Offset:   pageOffset * limit,
		Limit:    limit,
	})

	pageTotal := int(math.Ceil(float64(total) / float64(limit)))
	taskTotal := total

	projects, total, err := s.ProjectService.Find(c.Request().Context(), domain.ProjectFilter{
		Offset: 0,
		Limit:  100,
	})

	if total > 100 {
		c.Logger().Warnf("truncated projects list, total: %d, show: 100", total)
	}

	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))

		case domain.ENOTFOUND:
			return c.HTML(http.StatusNotFound, "no tasks found, unimplemented error handling")

		case domain.EINTERNAL:
			c.Logger().Error(err)
			return c.HTML(http.StatusInternalServerError, "internal error")
		}
	}

	return render(c, http.StatusOK, components.TaskListFind(
		generic.Map(taskItems, toTaskItemView),
		generic.Map(projects, toProjectView),
		pageOffset,
		pageTotal,
		taskTotal,
	))
}

// handleTaskDetailShow returns a task object.
// Data returned by handleTaskDetailShow is
// shown in Task Update screen
func (s *Server) handleTaskUpdateShow(c echo.Context) error {

	taskID := c.Param("id")
	if taskID == "" {
		return c.HTML(http.StatusBadRequest, "malformed url")
	}

	taskItem, err := s.TaskItemService.ByID(c.Request().Context(), taskID)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED, domain.ENOTFOUND:
			return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return c.HTML(http.StatusInternalServerError, "internal error")
		}
	}

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

	return render(c, http.StatusOK, pages.TaskUpdate(
		toTaskItemView(taskItem),
		[]string{"next", "school", "personal"},
		priorities,
		generic.Map(projects, toProjectView),
		"/logout",
	))
}

func (s *Server) handleTaskUpdate(c echo.Context) error {

	type formValue struct {
		Description string `form:"description" validate:"required"`
		Deadline    string `form:"deadline"`
		Schedule    string `form:"schedule"`
		Wait        string `form:"wait"`
		TagsJSON    string `form:"tags"`
		Priority    string `form:"priority"`
		ProjectID   string `form:"project_id"`
	}

	taskID := c.Param("id")
	if taskID == "" {
		return c.HTML(http.StatusBadRequest, "corrupted url")
	}

	form := formValue{}

	if err := c.Bind(&form); err != nil {
		return render(c, http.StatusBadRequest, components.AlertError("invalid form data"))
	}

	if err := c.Validate(form); err != nil {
		return render(c, http.StatusBadRequest, components.AlertError("invalid task info"))
	}

	tags, err := unmarshalTagsJSON([]byte(form.TagsJSON))
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError("invalid tags format"))
	}

	task, err := s.TaskService.Update(c.Request().Context(), taskID, domain.UpdateTaskCmd{
		Description: form.Description,
		Deadline:    form.Deadline,
		Schedule:    form.Schedule,
		Wait:        form.Wait,
		Priority:    form.Priority,
		Tags:        tags,
	})
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
		}
	}

	c.Logger().Info("update task successfully ", task.ID)

	c.Response().Header().Set("HX-Redirect", "/tasks")
	return c.NoContent(http.StatusFound)
}

func (s *Server) handleTaskSetProject(c echo.Context) error {

	taskID := c.Param("id")
	if taskID == "" {
		return c.HTML(http.StatusBadRequest, "malformed url")
	}

	type formValues struct {
		ProjectID string `form:"project-id"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		return c.HTML(http.StatusBadRequest, "invalid form data")
	}

	var projectID *string
	if form.ProjectID != "" {
		projectID = &form.ProjectID
	}

	task, err := s.TaskService.SetProject(c.Request().Context(), taskID, projectID)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return c.HTML(http.StatusInternalServerError, "internal error")
		}
	}

	c.Logger().Infof("set task project successful ID: %s projectID: %s", task.ID, task.ProjectID)

	// fetch task item to update the view
	taskItem, err := s.TaskItemService.ByID(c.Request().Context(), taskID)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
		}
	}

	return render(c, http.StatusOK, components.AlertAndUpdateTaskItemHidden(toTaskItemView(taskItem), "task project set successful"))
}

func (s *Server) handleTaskStart(c echo.Context) error {

	taskID := c.Param("id")
	if taskID == "" {
		return c.HTML(http.StatusBadRequest, "malformed url")
	}

	task, err := s.TaskService.Start(c.Request().Context(), taskID)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
		}
	}

	c.Logger().Infof("start task successful ID: %s", task.ID)

	// fetch task item to update the view
	taskItem, err := s.TaskItemService.ByID(c.Request().Context(), taskID)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
		}
	}

	return render(c, http.StatusOK, components.AlertAndUpdateTaskItemContent(
		toTaskItemView(taskItem),
		"task started successfully"),
	)
}

func (s *Server) handleTaskComplete(c echo.Context) error {

	taskID := c.Param("id")
	if taskID == "" {
		return render(c, http.StatusBadRequest, components.AlertError("malformed url"))

	}

	task, err := s.TaskService.Complete(c.Request().Context(), taskID)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return render(c, http.StatusInternalServerError, components.AlertError(("internal error")))
		}
	}

	c.Logger().Infof("task completed successful %s", task.ID)

	return render(c, http.StatusOK, components.AlertAndDeleteTaskItem("task completed successfully"))
}

func (s *Server) handleTaskDelete(c echo.Context) error {

	taskID := c.Param("id")
	if taskID == "" {
		return render(c, http.StatusBadRequest, components.AlertError("malformed url"))

	}

	task, err := s.TaskService.Delete(c.Request().Context(), taskID)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return render(c, http.StatusInternalServerError, components.AlertError(("internal error")))
		}
	}

	c.Logger().Infof("task deleted successful %s", task.ID)

	return render(c, http.StatusOK, components.AlertAndDeleteTaskItem("task deleted successfully"))
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

	tags, err := unmarshalTagsJSON([]byte(form.TagsJSON))
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
