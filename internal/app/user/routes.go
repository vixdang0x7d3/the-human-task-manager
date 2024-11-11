package user

import "github.com/labstack/echo/v4"

func (h *UserHandler) Route(e *echo.Group) {
	e.GET("/profile", h.HandleShowProfile)
	e.GET("/login", h.HandleShowLogin)
}
