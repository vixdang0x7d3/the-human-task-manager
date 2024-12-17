package http

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/generic"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/models"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/components"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

func (s *Server) registerProjectRoutes(r *echo.Group) {
	r.GET("/projects", s.handleProjectIndex)
	r.GET("/project/:id", s.handleProjectDetailShow)

	r.GET("/project/:id/tab-tasks", s.handleTabTasksShow)
	r.GET("/project/:id/tab-members", s.handleTabMembersShow)
	r.GET("/project/:id/tab-statistics", s.handleTabStatisticsShow)
	r.POST("/project/:id/tab-tasks/find", s.handleProjectTaskFind)

	r.POST("/invite-member/:project-id", s.handleInviteMember)
	r.GET("/accept-invitation/:project-id", s.handleAcceptInvitation)
	r.DELETE("/denine-invitation/:project-id", s.handleDenineInvitation)

	r.POST("/join-project-request", s.handleJoinProjectRequest)
	r.GET("/project/:id/accept-request/:user-id", s.handleAcceptRequest)
	r.DELETE("/project/:id/denine-request/:user-id", s.handleDenineRequest)

	r.DELETE("/leave-project/:project-id", s.handleLeaveProject)

	r.POST("/save-create-project", s.handleCreateProject)
	r.DELETE("/delete-project/:project-id", s.handleDeleteProject)
}

const (
	taskLimit       int = 2
	projectLimit    int = 2
	membershipLimit int = 2
)

func (s *Server) handleProjectIndex(c echo.Context) error {

	var pageOffset int
	if qParam := c.QueryParam("pageOffset"); qParam != "" {
		parsed, err := strconv.ParseInt(qParam, 10, 64)
		if err == nil {
			pageOffset = int(parsed)
		}
	}

	projectItem, total, err := s.ProjectmembershipService.Find(c.Request().Context(), domain.ProjectMembershipFilter{
		Offset: pageOffset * projectLimit,
		Limit:  projectLimit,
	})

	pageTotal := int(math.Ceil(float64(total) / float64(projectLimit)))
	projectTotal := total

	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(domain.ErrorMessage(err))
			return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
		}
	}

	trigger := c.Request().Header.Get("HX-Trigger")
	if trigger == "prev-btn" || trigger == "next-btn" {
		return render(c, http.StatusOK, components.ProjectList(
			generic.Map(projectItem, toProjectMemberShipItemView),
			pageOffset,
			pageTotal,
			projectTotal,
		))

	}

	return render(c, http.StatusOK, pages.ProjectIndex(
		generic.Map(projectItem, toProjectMemberShipItemView),
		pageOffset,
		pageTotal,
		projectTotal,
		"/logout",
	))
}

func toProjectItemView(projectItem domain.ProjectMembershipItem) models.ProjectMembershipItemView {
	return models.ProjectMembershipItemView{
		ProjectID: projectItem.Project.ID.String(),
		Title:     projectItem.Project.Title,
		Role:      projectItem.Role,
	}
}

func (s *Server) handleCreateProject(c echo.Context) error {
	type formValues struct {
		ProjectTitle string `form:"project-title" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusBadRequest, components.AlertError("invalid form data"))
	}

	if err := c.Validate(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid title")
	}

	//create project in db
	s.ProjectService.Create(c.Request().Context(), form.ProjectTitle)

	c.Response().Header().Set("HX-Redirect", "/projects")
	return c.NoContent(http.StatusOK)
}

func (s *Server) handleDeleteProject(c echo.Context) error {
	projectId := c.Param("project-id")
	if projectId == "" {
		return render(c, http.StatusBadRequest, components.AlertError("Request param error"))
	}

	_, err := s.ProjectService.Delete(c.Request().Context(), projectId)
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusInternalServerError, "internal error")
	}

	c.Response().Header().Set("HX-Redirect", "/projects")
	return c.NoContent(http.StatusOK)
}

func (s *Server) handleJoinProjectRequest(c echo.Context) error {
	type formValues struct {
		ProjectUUID string `form:"project-UUID" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusBadRequest, components.AlertError("invalid form data"))
	}

	if err := c.Validate(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid title")
	}

	_, err := s.ProjectmembershipService.Request(c.Request().Context(), domain.ProjectMembershipCmd{
		ProjectID: form.ProjectUUID,
	})

	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, err.Error())
	}

	c.Response().Header().Set("HX-Redirect", "/projects")
	return c.NoContent(http.StatusOK)
}

func (s *Server) handleAcceptRequest(c echo.Context) error {
	projectId := c.Param("id")
	if projectId == "" {
		return render(c, http.StatusBadRequest, components.AlertError("Request param error"))
	}

	userId := c.Param("user-id")
	if userId == "" {
		return render(c, http.StatusBadRequest, components.AlertError("Request param error"))
	}

	_, err := s.ProjectmembershipService.AcceptRequest(c.Request().Context(), domain.ProjectMembershipCmd{
		UserID:    userId,
		ProjectID: projectId,
	})
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError(err.Error()))
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/project/%s", projectId))
	return c.NoContent(http.StatusOK)
}

func (s *Server) handleDenineRequest(c echo.Context) error {
	projectId := c.Param("id")
	if projectId == "" {
		return render(c, http.StatusBadRequest, components.AlertError("Request param error"))
	}

	userId := c.Param("user-id")
	if userId == "" {
		return render(c, http.StatusBadRequest, components.AlertError("Request param error"))
	}

	_, err := s.ProjectmembershipService.Delete(c.Request().Context(), domain.ProjectMembershipCmd{
		UserID:    userId,
		ProjectID: projectId,
	})
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError(err.Error()))
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/project/%s", projectId))
	return c.NoContent(http.StatusOK)
}

func (s *Server) handleInviteMember(c echo.Context) error {
	projectId := c.Param("project-id")
	if projectId == "" {
		return render(c, http.StatusBadRequest, components.AlertError("Request param error"))
	}

	type formValues struct {
		MemberUUID string `form:"member-UUID" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusBadRequest, components.AlertError("invalid form data"))
	}

	if err := c.Validate(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid title")
	}

	_, err := s.ProjectmembershipService.Invite(c.Request().Context(), domain.ProjectMembershipCmd{
		ProjectID: projectId,
		UserID:    form.MemberUUID,
	})

	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, err.Error())
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/project/%s", projectId))
	return c.NoContent(http.StatusOK)
}

func (s *Server) handleAcceptInvitation(c echo.Context) error {
	projectId := c.Param("project-id")
	if projectId == "" {
		return render(c, http.StatusBadRequest, components.AlertError("Request param error"))
	}

	_, err := s.ProjectmembershipService.AcceptInvitation(c.Request().Context(), domain.ProjectMembershipCmd{
		ProjectID: projectId,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, err.Error())
	}

	c.Response().Header().Set("HX-Redirect", "/projects")
	return c.NoContent(http.StatusOK)
}

func (s *Server) handleDenineInvitation(c echo.Context) error {
	projectId := c.Param("project-id")
	if projectId == "" {
		return c.HTML(http.StatusBadRequest, "Request param error")
	}

	_, err := s.ProjectmembershipService.Delete(c.Request().Context(), domain.ProjectMembershipCmd{
		ProjectID: projectId,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, err.Error())
	}

	c.Response().Header().Set("HX-Redirect", "/projects")
	return c.NoContent(http.StatusOK)
}

func (s *Server) handleLeaveProject(c echo.Context) error {
	projectId := c.Param("project-id")
	if projectId == "" {
		return c.HTML(http.StatusBadRequest, "Request param error")
	}

	_, err := s.ProjectmembershipService.Delete(c.Request().Context(), domain.ProjectMembershipCmd{
		ProjectID: projectId,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, err.Error())
	}

	c.Response().Header().Set("HX-Redirect", "/projects")
	return c.NoContent(http.StatusOK)
}

func (s *Server) handleProjectDetailShow(c echo.Context) error {

	var pageOffset int
	if qParam := c.QueryParam("pageOffset"); qParam != "" {
		parsed, err := strconv.ParseInt(qParam, 10, 64)
		if err == nil {
			pageOffset = int(parsed)
		}
	}

	projectId := c.Param("id")
	if projectId == "" {
		return c.HTML(http.StatusBadRequest, "Request param error")
	}

	project, err := s.ProjectService.ByID(c.Request().Context(), projectId)
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
	}

	currentUserID := domain.UserIDFromContext(c.Request().Context())
	if currentUserID == nil {
		c.Logger().Error("no user in current session")
		return echo.NewHTTPError(http.StatusBadRequest, "no user in login session")
	}

	taskItems, total, err := s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		ProjectID: &projectId,
		Limit:     100,
	})

	if total > 100 {
		c.Logger().Warnf("truncate task list, total: %d show 100", total)
	}

	var (
		nCompleted  = 0
		nTotal      = 0
		percentDone = 0
		remain      = 0
	)

	for _, v := range taskItems {
		if v.State != domain.TaskStateDeleted {
			if v.State == domain.TaskStateCompleted {
				nCompleted += 1
			}
			nTotal += 1
		}
	}

	if nTotal != 0 {
		percentDone = int(nCompleted * 100 / nTotal)
		remain = nTotal - nCompleted
	}

	state := domain.TaskStateStarted
	taskItems, total, err = s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		ProjectID: &projectId,
		State:     &state,
		Offset:    pageOffset * taskLimit,
		Limit:     taskLimit,
	})

	pageTotal := int(math.Ceil(float64(total) / float64(taskLimit)))
	taskTotal := total

	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(domain.ErrorMessage(err))
			return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
		case domain.ENOTFOUND:
			return render(c, http.StatusOK, pages.ProjectDetail(
				true,
				toProjectView(project),
				generic.Map(taskItems, toTaskItemView),
				pageOffset,
				pageTotal,
				taskTotal,
				currentUserID.String(),
				"/logout",
				percentDone,
				remain,
			))
		}
	}

	trigger := c.Request().Header.Get("HX-Trigger")
	if trigger == "prev-btn" || trigger == "next-btn" {
		return render(c, http.StatusOK, components.TaskListProject(
			generic.Map(taskItems, toTaskItemView),
			toProjectView(project),
			pageOffset,
			pageTotal,
			taskTotal,
		))
	}

	return render(c, http.StatusOK, pages.ProjectDetail(
		false,
		toProjectView(project),
		generic.Map(taskItems, toTaskItemView),
		pageOffset,
		pageTotal,
		taskTotal,
		currentUserID.String(),
		"/logout",
		percentDone,
		remain,
	))
}

func (s *Server) handleTabTasksShow(c echo.Context) error {

	var pageOffset int
	if qParam := c.QueryParam("pageOffset"); qParam != "" {
		parsed, err := strconv.ParseInt(qParam, 10, 64)
		if err == nil {
			pageOffset = int(parsed)
		}
	}

	projectId := c.Param("id")
	if projectId == "" {
		return c.HTML(http.StatusBadRequest, "Request param error")
	}

	project, err := s.ProjectService.ByID(c.Request().Context(), projectId)
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
	}

	taskItems, total, err := s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		ProjectID: &projectId,
		Limit:     100,
	})

	if total > 100 {
		c.Logger().Warnf("truncated task list, total: %d show: 100")
	}

	var (
		nCompleted  = 0
		nTotal      = 0
		percentDone = 0
		remain      = 0
	)
	for _, v := range taskItems {
		if v.State != domain.TaskStateDeleted {
			if v.State == domain.TaskStateCompleted {
				nCompleted += 1
			}
			nTotal += 1
		}
	}

	if nTotal != 0 {
		percentDone = int(nCompleted * 100 / nTotal)
		remain = nTotal - nCompleted
	}

	state := domain.TaskStateStarted
	taskItems, total, err = s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		ProjectID: &projectId,
		State:     &state,
		Offset:    pageOffset * taskLimit,
		Limit:     taskLimit,
	})

	pageTotal := int(math.Ceil(float64(total) / float64(taskLimit)))
	taskTotal := total

	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(domain.ErrorMessage(err))
			return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
		case domain.ENOTFOUND:
			return render(c, http.StatusOK, components.ProjectTasks(
				true,
				toProjectView(project),
				generic.Map(taskItems, toTaskItemView),
				pageOffset,
				pageTotal,
				taskTotal,
				percentDone,
				remain,
			))
		}
	}

	trigger := c.Request().Header.Get("HX-Trigger")
	if trigger == "prev-btn" || trigger == "next-btn" {
		return render(c, http.StatusOK, components.TaskListProject(
			generic.Map(taskItems, toTaskItemView),
			toProjectView(project),
			pageOffset,
			pageTotal,
			taskTotal,
		))

	}

	return render(c, http.StatusOK, components.ProjectTasks(
		false,
		toProjectView(project),
		generic.Map(taskItems, toTaskItemView),
		pageOffset,
		pageTotal,
		taskTotal,
		percentDone,
		remain,
	))
}

func (s *Server) handleTabMembersShow(c echo.Context) error {

	projectId := c.Param("id")
	if projectId == "" {
		return c.HTML(http.StatusBadRequest, "Request param error")
	}

	project, err := s.ProjectService.ByID(c.Request().Context(), projectId)
	if err != nil {
		return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
	}

	var pageOffset int
	if qParam := c.QueryParam("pageOffset"); qParam != "" {
		parsed, err := strconv.ParseInt(qParam, 10, 64)
		if err == nil {
			pageOffset = int(parsed)
		}
	}

	memberItems, total, err := s.ProjectmembershipService.Find(c.Request().Context(), domain.ProjectMembershipFilter{
		ProjectID: &projectId,
		Offset:    pageOffset * membershipLimit,
		Limit:     membershipLimit,
	})

	currentUserID := domain.UserIDFromContext(c.Request().Context())
	if currentUserID == nil {
		c.Logger().Error("no user in current session")
		return echo.NewHTTPError(http.StatusBadRequest, "no user in login session")
	}

	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
		}
	}

	pageTotal := int(math.Ceil(float64(total) / float64(membershipLimit)))
	memberTotal := total

	trigger := c.Request().Header.Get("HX-Trigger")
	if trigger == "prev-btn" || trigger == "next-btn" {
		return render(c, http.StatusOK, components.MemberList(
			toProjectView(project),
			generic.Map(memberItems, toMemberItemView),
			pageOffset,
			pageTotal,
			memberTotal,
			currentUserID.String(),
		))

	}

	return render(c, http.StatusOK, components.ProjectMembers(
		toProjectView(project),
		generic.Map(memberItems, toMemberItemView),
		pageOffset,
		pageTotal,
		memberTotal,
		currentUserID.String(),
	))
}

func toProjectMemberShipItemView(membership domain.ProjectMembershipItem) models.ProjectMembershipItemView {
	return models.ProjectMembershipItemView{
		UserID:    membership.User.ID.String(),
		ProjectID: membership.Project.ID.String(),
		Title:     membership.Project.Title,
		Username:  membership.User.Username,
		Role:      membership.Role,
	}
}

func toMemberItemView(memberItem domain.ProjectMembershipItem) models.ProjectMembershipItemView {
	return models.ProjectMembershipItemView{
		UserID:    memberItem.User.ID.String(),
		ProjectID: memberItem.Project.ID.String(),
		Username:  memberItem.User.Username,
		Role:      memberItem.Role,
	}
}

func (s *Server) handleTabStatisticsShow(c echo.Context) error {

	projectId := c.Param("id")
	if projectId == "" {
		return render(c, http.StatusBadRequest, components.AlertError("Request param error"))
	}

	project, err := s.ProjectService.ByID(c.Request().Context(), projectId)
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
	}

	prj := models.ProjectView{
		Title: project.Title,
		ID:    projectId,
	}

	members, totalMember, err := s.ProjectmembershipService.Find(c.Request().Context(), domain.ProjectMembershipFilter{
		ProjectID: &projectId,
		Limit:     100,
		Offset:    0,
	})
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return render(c, http.StatusBadRequest, components.AlertError("internal error"))
		}
	}
	if totalMember > 100 {
		c.Logger().Warnf("truncated project members, total: %d, show: 100", totalMember)
	}

	memberMaps := map[string]int{}

	for _, mem := range members {
		memberMaps[mem.User.Username] = 0
	}

	taskItems, totalTask, err := s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		ProjectID: &projectId,
		Limit:     100,
		Offset:    0,
	})
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError(err.Error()))
	}
	if totalTask > 100 {
		c.Logger().Warnf("truncated project members, total: %d, show: 100", totalTask)
	}

	for _, item := range taskItems {
		if item.CompletedByName != "" {
			memberMaps[item.CompletedByName] += 1
		}
	}

	Usernames := []string{}
	CompleteNums := []int{}

	for k, v := range memberMaps {
		Usernames = append(Usernames, k)
		CompleteNums = append(CompleteNums, v)
	}

	var (
		nCompleted  = 0
		nTotal      = 0
		percentDone = 0
		remain      = 0
	)
	for _, v := range taskItems {
		if v.State != domain.TaskStateDeleted {
			if v.State == domain.TaskStateCompleted {
				nCompleted += 1
			}
			nTotal += 1
		}
	}

	if nTotal != 0 {
		percentDone = int(nCompleted * 100 / nTotal)
		remain = nTotal - nCompleted
	}

	return render(c, http.StatusOK, components.ProjectStatistics(
		Usernames,
		CompleteNums,
		prj,
		totalTask,
		percentDone,
		remain,
	))
}

func (s *Server) handleProjectTaskFind(c echo.Context) error {

	type formValues struct {
		Query    string `form:"query"`
		Priority string `form:"priority"`
		State    string `form:"state"`
		Days     string `form:"days"`
		Months   string `form:"months"`
	}

	projectID := c.Param("id")
	if projectID == "" {
		return render(c, http.StatusBadRequest, components.AlertError("malformed url"))
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		return render(c, http.StatusBadRequest, components.AlertError("invalid form data"))
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
	}

	if form.Months != "" {
		val, err := strconv.ParseInt(form.Months, 10, 64)
		if err == nil && val != 0 {
			months = &val

		}
	}

	var pageOffset int
	if qParam := c.QueryParam("pageOffset"); qParam != "" {
		parsed, err := strconv.ParseInt(qParam, 10, 64)
		if err == nil {
			pageOffset = int(parsed)
		}
	}

	taskItems, total, err := s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		ProjectID: &projectID,
		Q:         query,
		Priority:  priority,
		State:     state,
		Days:      days,
		Months:    months,
		Offset:    pageOffset * taskLimit,
		Limit:     taskLimit,
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

	pageTotal := int(math.Ceil(float64(total) / float64(taskLimit)))
	taskTotal := total

	project := models.ProjectView{
		ID: projectID,
	}

	return render(c, http.StatusOK, components.TaskListProjectFind(
		generic.Map(taskItems, toTaskItemView),
		project,
		pageOffset,
		pageTotal,
		taskTotal,
	))
}
