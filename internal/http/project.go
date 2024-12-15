package http

import (
	"fmt"
	"net/http"

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

func (s *Server) handleProjectIndex(c echo.Context) error {
	projectItem, total, err := s.ProjectmembershipService.Find(c.Request().Context(), domain.ProjectMembershipFilter{
		Offset: 0,
		Limit:  10,
	})

	// handle pagination
	_ = total

	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(domain.ErrorMessage(err))
			return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
		}
	}

	return render(c, http.StatusOK, pages.ProjectIndex(
		generic.Map(projectItem, toProjectItemView),
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

	projectId := c.Param("id")
	if projectId == "" {
		return c.HTML(http.StatusBadRequest, "Request param error")
	}

	project, err := s.ProjectService.ByID(c.Request().Context(), projectId)
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
	}

	prj := models.ProjectView{
		UserID: project.UserID.String(),
		Title:  project.Title,
		ID:     projectId,
	}

	currentUserID := domain.UserIDFromContext(c.Request().Context())
	if currentUserID == nil {
		c.Logger().Error("no user in current session")
		return echo.NewHTTPError(http.StatusBadRequest, "no user in login session")
	}

	state := "started"
	taskItems, total, err := s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		ProjectID: &projectId,
		State:     &state,
		Offset:    0,
		Limit:     100,
	})

	// handle pagination
	_ = total

	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(domain.ErrorMessage(err))
			return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
		case domain.ENOTFOUND:
			return render(c, http.StatusOK, pages.ProjectDetail(true,
				prj,
				generic.Map(taskItems, toTaskItemView),
				currentUserID.String(),
				"/logout"))
		}
	}

	return render(c, http.StatusOK, pages.ProjectDetail(false, prj,
		generic.Map(taskItems, toTaskItemView),
		currentUserID.String(),
		"/logout"))
}

func (s *Server) handleTabTasksShow(c echo.Context) error {

	projectId := c.Param("id")
	if projectId == "" {
		return c.HTML(http.StatusBadRequest, "Request param error")
	}

	project, err := s.ProjectService.ByID(c.Request().Context(), projectId)
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
	}

	prj := models.ProjectView{
		Title: project.Title,
		ID:    projectId,
	}

	state := "started"
	taskItems, total, err := s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		ProjectID: &projectId,
		State:     &state,
		Offset:    0,
		Limit:     100,
	})

	// handle pagination
	_ = total

	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
		case domain.EINTERNAL:
			c.Logger().Error(domain.ErrorMessage(err))
			return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
		case domain.ENOTFOUND:
			return render(c, http.StatusOK, components.ProjectTasks(true, prj,
				generic.Map(taskItems, toTaskItemView)))
		}
	}

	return render(c, http.StatusOK, components.ProjectTasks(false, prj,
		generic.Map(taskItems, toTaskItemView)))
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

	prj := models.ProjectView{
		Title: project.Title,
		ID:    projectId,
	}

	memberiItems, total, err := s.ProjectmembershipService.Find(c.Request().Context(), domain.ProjectMembershipFilter{
		ProjectID: &projectId,
		Offset:    0,
		Limit:     15,
	})

	_ = total

	return render(c, http.StatusOK, components.ProjectMembers(prj,
		generic.Map(memberiItems, toMemberItemView)))
}

func toMemberItemView(memberiItem domain.ProjectMembershipItem) models.ProjectMembershipItemView {
	return models.ProjectMembershipItemView{
		UserID:    memberiItem.User.ID.String(),
		ProjectID: memberiItem.Project.ID.String(),
		Username:  memberiItem.User.Username,
		Role:      memberiItem.Role,
	}
}

func (s *Server) handleTabStatisticsShow(c echo.Context) error {

	projectId := c.Param("id")
	if projectId == "" {
		return c.HTML(http.StatusBadRequest, "Request param error")
	}

	project, err := s.ProjectService.ByID(c.Request().Context(), projectId)
	if err != nil {
		return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
	}

	prj := models.ProjectView{
		Title: project.Title,
		ID:    projectId,
	}

	return render(c, http.StatusOK, components.ProjectStatistics(prj))
}
