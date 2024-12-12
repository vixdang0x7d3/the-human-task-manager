package http_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	htmhttp "github.com/vixdang0x7d3/the-human-task-manager/internal/http"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/mock"
)

type Server struct {
	*htmhttp.Server

	UserService mock.UserService
	TaskService mock.TaskService
}

func MustOpenServer(tb testing.TB) *Server {

	s := &Server{Server: htmhttp.NewServer(logrus.New())}

	s.Server.UserService = &s.UserService
	s.Server.UserService = &s.UserService

	if err := s.Server.Open(); err != nil {
		tb.Fatal(err)
	}

	return s
}

func CloseServer(tb testing.TB, s *Server) {
	tb.Helper()
	s.Server.Close()
}

func (s *Server) MustNewRequest(tb testing.TB, ctx context.Context, method, url string, body io.Reader) *http.Request {
	tb.Helper()
	r := httptest.NewRequest(method, s.URL()+url, body)

	return r
}
