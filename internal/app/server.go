package app

import (
	"io/fs"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	userapp "github.com/vixdang0x7d3/the-human-task-manager/internal/app/user"
	usercore "github.com/vixdang0x7d3/the-human-task-manager/internal/core/user"

	"github.com/vixdang0x7d3/the-human-task-manager/internal/app/validate"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/template"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/template/pages"
)

func NewServer(db *database.Queries, staticAssets fs.FS) (*echo.Echo, error) {
	e := echo.New()

	e.Validator = &validate.CustomValidator{Validator: validator.New()}
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
		Browse:     true,
		Filesystem: http.FS(staticAssets),
	}))
	route(e, db)

	return e, nil
}

func route(e *echo.Echo, db *database.Queries) {

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

	userService := usercore.New(db)
	userHandler := userapp.NewHandler(userService)
	userHandler.Route(v1)
}
