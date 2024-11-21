package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
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
		fmt.Sprintf(":%s", os.Getenv("ADDRESS")),
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
	server *http.Server
	db     *postgres.DB
}

func NewMain(addr, dburl string) *Main {
	return &Main{
		Address: addr,
		db:      postgres.NewDB(dburl),
		server:  http.NewServer(),
	}
}

func (m *Main) Run(ctx context.Context) (err error) {
	if err = m.db.Open(ctx); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	m.server.Addr = m.Address
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
