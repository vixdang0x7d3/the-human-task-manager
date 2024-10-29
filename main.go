package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vixdang0x7d3/the-human-task-manager/internal"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/app"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
)

func main() {
	godotenv.Load()
	servePort := os.Getenv("PORT")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	defer conn.Close(context.Background())

	db := database.New(conn)

	c := domain.UserCore{Store: db}
	userHandler := app.UserHandler{
		Service: &c,
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{
"time":"${time_rfc3339_nano}"
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
"bytes_out":${bytes_out}
}\n`,
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
	v1.GET("/err", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "response with error always")
	})

	v1.POST("/users", userHandler.HandleCreateUser)
	v1.GET("/users/:id", userHandler.HandleGetUser)

	e.Logger.Fatal(e.Start(":" + servePort))
}
