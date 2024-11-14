package app

import "github.com/labstack/echo/v4"

func (h *UserHandler) Route(e *echo.Group) {

	// ui serving endpoints
	e.GET("/profile", h.HandleShowProfile)

	e.GET("/home", h.HandleShowHome)

	e.GET("/login", h.HandleShowLogin)
	e.GET("/signup", h.HandleShowSignup)

	// api endpoints
	e.POST("/users", h.HandleUserCreate)
	e.POST("/login-email", h.HandleLoginCheckEmail)
	e.POST("/login", h.HandleLoginCheckPassword)
}
