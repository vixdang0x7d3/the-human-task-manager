package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/assets"

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

func NewServer(logger *logrus.Logger) *Server {
	s := &Server{
		echo:     echo.New(),
		sessions: scs.New(),
	}

	customLogger := customLogger{Logger: logger}

	s.echo.Logger = customLogger

	s.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_custom}] ${status} ${method} ${uri} ${latency_human}` + "\n" +
			`    Remote IP: ${remote_ip}` + "\n" +
			`    Host: ${host}` + "\n" +
			`    User Agent: ${user_agent}` + "\n" +
			`    Error: "${error}"` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
		Output:           os.Stdout,
	}))

	s.echo.Validator = &customValidator{Validator: validator.New()}

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
		return c.Redirect(http.StatusFound, "/index")
	})

	// registers unauthenticated routes
	{
		r := s.echo.Group("u", s.requireNoAuth(s.sessions))
		s.registerAuthRoutes(r)
	}

	// registers authenticated routes
	{
		r := s.echo.Group("", s.requireAuth(s.sessions))
		s.registerUserRoutes(r)
		s.registerTaskRoutes(r)
	}

	// registers authenticated routes
	{
		r := s.echo.Group("", s.requireAuth(s.sessions))
		s.registerCalendarRoutes(r)
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

func (s *Server) requireNoAuth(sessions *scs.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := sessions.GetString(c.Request().Context(), "userID")
			if userID != "" {

				user, err := s.UserService.ByID(c.Request().Context(), userID)
				if err != nil {
					c.Logger().Error("cannot find session user", userID, err)
				} else {
					r := c.Request().WithContext(domain.NewContextWithUser(c.Request().Context(), &user))
					c.SetRequest(r)
				}

				// using redirect header so htmx won't
				// inject the whole returned page into the
				// message div :D
				// FIX: use hx-boost & hx-push-url instead
				// then redirect normally
				c.Response().Header().Set("HX-Redirect", "../index")
				return c.NoContent(http.StatusTemporaryRedirect)
			}
			return next(c)
		}
	}
}

// TODO: implement url memorization
func (s *Server) requireAuth(sessions *scs.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := sessions.GetString(c.Request().Context(), "userID")
			if userID == "" {
				return c.Redirect(http.StatusTemporaryRedirect, "/u/login")
			}

			user, err := s.UserService.ByID(c.Request().Context(), userID)
			if err != nil {
				c.Logger().Error("cannot find session user", userID, err)
			} else {
				r := c.Request().WithContext(domain.NewContextWithUser(c.Request().Context(), &user))
				c.SetRequest(r)
			}

			return next(c)
		}
	}
}

// FIX: deprecated
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

type customLogger struct {
	*logrus.Logger
}

// Level returns logger level
func (l customLogger) Level() log.Lvl {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.InfoLevel:
		return log.INFO
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	default:
		return log.OFF
	}
}

func (l customLogger) SetHeader(_ string) {}

func (l customLogger) SetPrefix(_ string) {}

func (l customLogger) Prefix() string {
	return ""
}

func (l customLogger) SetLevel(lvl log.Lvl) {
	switch lvl {
	case log.DEBUG:
		l.Logger.SetLevel(logrus.DebugLevel)
	case log.INFO:
		l.Logger.SetLevel(logrus.InfoLevel)
	case log.WARN:
		l.Logger.SetLevel(logrus.WarnLevel)
	case log.ERROR:
		l.Logger.SetLevel(logrus.ErrorLevel)
	}
}

func (l customLogger) Output() io.Writer {
	return l.Logger.Out
}

func (l customLogger) SetOutput(w io.Writer) {
	l.Logger.SetOutput(w)
}

func (l customLogger) Printj(j log.JSON) {}

func (l customLogger) Debugj(j log.JSON) {}

func (l customLogger) Errorj(j log.JSON) {}

func (l customLogger) Fatalj(j log.JSON) {}

func (l customLogger) Infoj(j log.JSON) {}

func (l customLogger) Panicj(j log.JSON) { panic("not implemented") }

func (l customLogger) Warnj(j log.JSON) {}
