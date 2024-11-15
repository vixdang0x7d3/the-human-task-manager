package sdk

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
)

func RequireAuth(sessionManager *scs.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := sessionManager.GetString(c.Request().Context(), "userID")
			if userID == "" {
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}

			return next(c)
		}
	}
}
