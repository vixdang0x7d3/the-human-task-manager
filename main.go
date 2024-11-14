package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/app"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/database"
)

//go:embed static
var staticAssets embed.FS
var sessionManager *scs.SessionManager

func main() {
	godotenv.Load()
	servePort := fmt.Sprintf(":" + os.Getenv("PORT"))

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())
	db := database.New(conn)

	sessionManager = scs.New()

	e := app.SetupServer(db, staticAssets, sessionManager)
	e.Logger.Fatal(e.Start(servePort))
}
