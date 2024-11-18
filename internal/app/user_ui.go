package app

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/template"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/template/pages"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/types"
)

func (h *UserHandler) HandleShowProfile(c echo.Context) error {

	userIDStr := h.SessionManager.GetString(c.Request().Context(), "userID")
	if userIDStr == "" {
		return c.HTML(http.StatusUnauthorized, "Unauthorized user")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.HTML(http.StatusBadRequest, "Invalid session data")
	}

	user, err := h.Service.ByID(c.Request().Context(), userID)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "Internal Server Error")
	}

	viewData := types.UserViewModel{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return template.Render(c, http.StatusOK, pages.Profile(viewData))
}

func (h *UserHandler) HandleShowSignup(c echo.Context) error {

	return template.Render(c, http.StatusOK, pages.Signup())
}

func (h *UserHandler) HandleShowLogin(c echo.Context) error {

	return template.Render(c, http.StatusOK, pages.Login())
}

func (h *UserHandler) HandleShowHome(c echo.Context) error {

	userIDStr := h.SessionManager.GetString(c.Request().Context(), "userID")
	if userIDStr == "" {
		return c.HTML(http.StatusUnauthorized, "Unauthorized user")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.HTML(http.StatusBadRequest, "Invalid session data")
	}

	user, err := h.Service.ByID(c.Request().Context(), userID)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "Internal Server Error")
	}

	viewData := types.UserViewModel{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return template.Render(c, http.StatusOK, pages.Home(viewData))
}

func (h *UserHandler) HandleShowTaskList(c echo.Context) error {

	return template.Render(c, http.StatusOK, pages.Tasklist())
}

func (h *UserHandler) HandleShowTaskDetail(c echo.Context) error {

	taskdata := types.TaskViewModel{
		Title:       "This is a Title of the Task",
		Description: "This is a Decription of the Task",
		Create_at:   "11:57 17/11/2024",
		Update_at:   "17/11/2024 11:57",
		Deadline:    "2024-11-25T23:59",
		Schedule:    "none",
		Tags:        []string{"school", "math", "ttcs"},
	}

	Tags := []string{"morning", "personal", "groceries"}

	return template.Render(c, http.StatusOK, pages.TaskDetail(taskdata, Tags))
}
