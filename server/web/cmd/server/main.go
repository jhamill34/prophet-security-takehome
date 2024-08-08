package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/jhamill34/prophet-security-takehome/server/database/pkg/database"
)

const (
	PORT = 3333
)

func main() {
	db := NewDatabase(context.TODO(), "host=localhost port=5432 user=prophet-th password=prophet-th dbname=prophet-th sslmode=disable")
	queries := database.New(db)

	server := NewHttpServer(
		PORT,
		RootRouter(
			NewNodeResource(queries),
			NewSourceResource(queries),
			NewAllowListResource(queries),
		),
	)

	Run(server)
}

func Run(server *http.Server) {
	slog.Info("Starting server")
	server.ListenAndServe()
}
