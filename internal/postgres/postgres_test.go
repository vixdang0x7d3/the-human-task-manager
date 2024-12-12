package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/postgres"
	"github.com/vixdang0x7d3/the-human-task-manager/test_data/migrations"
)

var dsn string

func TestDB(t *testing.T) {
	db := MustOpenDB(t, context.Background())
	CloseDB(t, db)
}

func MustOpenDB(tb testing.TB, ctx context.Context) *postgres.DB {
	tb.Helper()

	db := postgres.NewDB(dsn)
	if err := db.Open(ctx); err != nil {
		tb.Fatal(err)
	}
	return db
}

func CloseDB(tb testing.TB, db *postgres.DB) {
	tb.Helper()
	db.Close()
}

func TestMain(m *testing.M) {

	var db *sql.DB

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not construct dockertest pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=htm_postgres",
			"listen_addresses='*'",
		},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	dsn = fmt.Sprintf("postgres://postgres:postgres@%s/htm_postgres?sslmode=disable", hostAndPort)

	log.Println("connecting to database on dsn: ", dsn)

	resource.Expire(120)

	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("pgx", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}
	mustMigrate(db)

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("could not purge resource: %s", err)
		}
	}()

	m.Run() // run all tests in this package
}

func mustMigrate(db *sql.DB) {
	goose.SetBaseFS(migrations.EmbedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "."); err != nil {
		log.Fatal(err)
	}
}
