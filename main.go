package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	godotenv.Load()
	m := NewMain(
		os.Getenv("ADDRESS"),
		os.Getenv("DB_URL"),
	)

	if err := m.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	<-ctx.Done()
	m.Close()
}

type Main struct {
	Address string

	// attaching to Main in order to do clean up in Close()
	logger *logrus.Logger
	server *http.Server
	db     *postgres.DB
}

func NewMain(addr, dsn string) *Main {

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
		FullTimestamp:   true,
		PadLevelText:    true,
	})

	return &Main{
		Address: addr,
		db:      postgres.NewDB(dsn),
		logger:  logger,
		server:  http.NewServer(logger),
	}
}

func (m *Main) Run(ctx context.Context) (err error) {
	if err = m.db.Open(ctx); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	userService := postgres.NewUserService(m.db, m.logger)
	taskService := postgres.NewTaskService(m.db, m.logger)
	projectService := postgres.NewProjectService(m.db, m.logger)
	taskItemService := postgres.NewTaskItemService(m.db, m.logger)

	m.server.Addr = m.Address
	m.server.UserService = userService
	m.server.TaskService = taskService
	m.server.ProjectService = projectService
	m.server.TaskItemService = taskItemService

	return m.server.Open()
}

func (m *Main) Close() {
	if m.server != nil {
		m.server.Close()
	}
	if m.db != nil {
		m.db.Close()
	}
}
