package app

import (
	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/app/sdk"
)

func (h *UserHandler) Route(e *echo.Echo) {

	// ui serving endpoints
	e.GET("/profile", h.HandleShowProfile, sdk.RequireAuth(h.SessionManager))
	e.GET("/home", h.HandleShowHome, sdk.RequireAuth((h.SessionManager)))
	e.GET("/tasklist", h.HandleShowTaskList)
	e.GET("/taskdetail", h.HandleShowTaskDetail)
	e.GET("/login", h.HandleShowLogin)
	e.GET("/signup", h.HandleShowSignup)

	// api endpoints
	e.POST("/users", h.HandleUserCreate)
	e.POST("/login-email", h.HandleLoginCheckEmail)
	e.POST("/login", h.HandleLoginCheckPassword)
	e.GET("/login-success-redirect", h.HandleLoginRedirect)
	e.GET("/logout", h.HandleLogout)
}
