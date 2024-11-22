package http

import "github.com/labstack/echo/v4"

func (s *Server) registerUserRoutes(r *echo.Group) {
	r.GET("/login", s.handleLoginShow)
	r.GET("/signup", s.handleSignupShow)

	r.POST("/signup", s.handleSignup)
	r.POST("/login-email", s.handleLoginEmail)
	r.POST("/login", s.handleLogin)
	r.GET("/login-success", s.handleLoginSuccess)
	r.DELETE("/logout", s.handleLogout)
}

func (s *Server) handleLoginShow(c echo.Context) error {
	return nil
}

func (s *Server) handleSignupShow(c echo.Context) error {
	return nil
}

func (s *Server) handleSignup(c echo.Context) error {
	return nil
}

func (s *Server) handleLoginEmail(c echo.Context) error {
	return nil
}

func (s *Server) handleLogin(c echo.Context) error {
	return nil
}

func (s *Server) handleLoginSuccess(c echo.Context) error {
	return nil
}

func (s *Server) handleLogout(c echo.Context) error {
	return nil
}
