package http

import "github.com/labstack/echo/v4"

func (s *Server) registerTaskRoutes(r *echo.Group) {
	r.GET("/tasks", s.handleTaskShow)
	r.GET("/tasks/create", s.handleTaskCreateShow)

	r.POST("/tasks/create", s.handleTaskCreate)
}

func (s *Server) handleTaskShow(c echo.Context) error {
	return nil
}

func (s *Server) handleTaskCreateShow(c echo.Context) error {
	return nil
}

func (s *Server) handleTaskCreate(c echo.Context) error {
	return nil
}
