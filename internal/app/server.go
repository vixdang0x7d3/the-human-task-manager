package app

import (
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	slogecho "github.com/samber/slog-echo"
	session "github.com/spazzymoto/echo-scs-session"

	"github.com/vixdang0x7d3/the-human-task-manager/internal/app/sdk"

	"github.com/vixdang0x7d3/the-human-task-manager/internal/core"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
)

func SetupServer(db *database.Queries, staticAssets fs.FS, sessionManager *scs.SessionManager) *echo.Echo {
	e := echo.New()

	// validator & middleware
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	e.Use(slogecho.New(logger))
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
	//	e.Use(sdk.RequireAuth(sessionManager))
	e.Validator = sdk.NewCustomValidator()

	// services
	userService := core.NewUserCore(db)

	// handlers
	userHandler := NewUserHandler(userService, sessionManager)

	// routes
	routeDefaults(e)
	userHandler.Route(e)

	return e
}
