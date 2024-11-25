package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/assets"

	slogecho "github.com/samber/slog-echo"
	session "github.com/spazzymoto/echo-scs-session"
)

// TODO: should have a server struct

type Server struct {
	echo     *echo.Echo
	sessions *scs.SessionManager

	Addr        string
	UserService domain.UserService
	TaskService domain.TaskService
}

func NewServer() *Server {
	s := &Server{
		echo:     echo.New(),
		sessions: scs.New(),
	}

	s.echo.Validator = &customValidator{Validator: validator.New()}

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

	s.echo.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "/",
		Filesystem: http.FS(assets.FS),
		Browse:     false,
	}))

	s.echo.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/home")
	})

	// registers unauthenticated routes
	{
		r := s.echo.Group("u", requireNoAuth(s.sessions))
		s.registerAuthRoutes(r)
	}

	// registers authenticated routes
	{
		r := s.echo.Group("", requireAuth(s.sessions))
		s.registerUserRoutes(r)
	}

	return s
}

func (s *Server) Open() error {
	go func() {
		if err := s.echo.Start(fmt.Sprintf(":%s", s.Addr)); err != nil && err != http.ErrServerClosed {
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

func (s *Server) URL() string {
	domain := "localhost"
	scheme := "http"
	return fmt.Sprintf("%s://%s:%s", scheme, domain, s.Addr)
}

func render(ctx echo.Context, code int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(code, buf.String())
}

func requireNoAuth(sessions *scs.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := sessions.GetString(c.Request().Context(), "userID")
			if userID != "" {
				c.Response().Header().Set("HX-Redirect", "../home")
				return c.NoContent(http.StatusOK)
			}
			return next(c)
		}
	}
}

// TODO: implement url memorization
func requireAuth(sessions *scs.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := sessions.GetString(c.Request().Context(), "userID")
			if userID == "" {
				return c.Redirect(http.StatusTemporaryRedirect, "/u/login")
			}
			return next(c)
		}
	}
}

func (s *Server) userIDFromSession(ctx context.Context) (uuid.UUID, error) {
	idString := s.sessions.GetString(ctx, "userID")
	if idString == "" {
		return uuid.Nil, domain.Errorf(domain.EUNAUTHORIZED, "unauthorized user")
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, domain.Errorf(domain.EUNAUTHORIZED, "invalid session data")
	}

	return id, nil
}

type customValidator struct {
	Validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}
