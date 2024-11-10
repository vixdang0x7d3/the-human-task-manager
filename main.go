package main

import (
	"context"
	"embed"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/app"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
)

//go:embed static
var staticAssets embed.FS

func main() {
	godotenv.Load()
	servePort := os.Getenv("PORT")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	defer conn.Close(context.Background())

	db := database.New(conn)

	e, err := app.NewServer(db, staticAssets)
	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":" + servePort))
}
