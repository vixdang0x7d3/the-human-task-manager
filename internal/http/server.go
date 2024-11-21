package http

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	slogecho "github.com/samber/slog-echo"
	session "github.com/spazzymoto/echo-scs-session"
)

// TODO: should have a server struct

type Server struct {
	echo     *echo.Echo
	sessions *scs.SessionManager

	Addr string
}

func NewServer() *Server {
	s := &Server{
		echo:     echo.New(),
		sessions: scs.New(),
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s.echo.Use(slogecho.New(logger))

	s.echo.Use(middleware.Recover())

	s.echo.Use(session.LoadAndSave(s.sessions))

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://*", "https://*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.echo.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "hi mom")
	})

	return s
}

func (s *Server) Open() error {
	go func() {
		if err := s.echo.Start(s.Addr); err != nil && err != http.ErrServerClosed {
			s.echo.Logger.Fatal("shutting down the server")
		}
	}()
	return nil
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.echo.Shutdown(ctx); err != nil {
		s.echo.Logger.Fatal(err)
	}
}
