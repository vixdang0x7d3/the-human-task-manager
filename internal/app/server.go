package app

import (
	"io/fs"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	session "github.com/spazzymoto/echo-scs-session"

	"github.com/vixdang0x7d3/the-human-task-manager/internal/app/sdk"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/template/pages"

	"github.com/vixdang0x7d3/the-human-task-manager/internal/core"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/template"
)

func SetupServer(db *database.Queries, staticAssets fs.FS, sessionManager *scs.SessionManager) *echo.Echo {
	e := echo.New()

	// validator & middleware
	e.Validator = sdk.NewCustomValidator()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://*", "https://*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "static",
		Filesystem: http.FS(staticAssets),
	}))
	e.Use(session.LoadAndSave(sessionManager))

	// services
	userService := core.NewUserCore(db)

	// handlers
	userHandler := NewUserHandler(userService, sessionManager)

	// routes
	v1 := e.Group("/v1")
	v1.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{}{})
	})
	v1.GET("/err", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "response with error always")
	})
	v1.GET("/", func(c echo.Context) error {
		return template.Render(c, http.StatusOK, pages.Index("ma"))
	})

	userHandler.Route(v1)

	return e
}
