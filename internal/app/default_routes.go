package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func routeDefaults(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		// return template.Render(c, http.StatusOK, pages.Index("ma"))
		return c.Redirect(http.StatusFound, "/home")
	})
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{}{})
	})
	e.GET("/err", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "response with error always")
	})

	// used to delete temporary htmx components :D
	e.GET("/empty", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "")
	})
}
