package main

import (
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vixdang0x7d3/the-human-task-manager/internal"
)

func main() {
	godotenv.Load()
	servePort := os.Getenv("PORT")

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}"
,"id":"${id}",
"remote_ip":"${remote_ip}",` +
			`"host":"${host}",
"method":"${method}",
"uri":"${uri}",
"user_agent":"${user_agent}",` +
			`"status":${status},
			"error":"${error}",
"latency":${latency},
"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},
"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))

	e.Validator = &internal.CustomValidator{Validator: validator.New()}
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://*", "https://*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1 := e.Group("/v1")

	v1.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{}{})
	})
	e.GET("/err", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong")
	})

	e.Logger.Fatal(e.Start(":" + servePort))
}
