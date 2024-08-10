package main

import (
	"context"
	"log/slog"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jhamill34/prophet-security-takehome/server/api/pkg/api"
	"github.com/jhamill34/prophet-security-takehome/server/database/pkg/database"
	"github.com/jhamill34/prophet-security-takehome/server/web/internal/auth"
	"github.com/jhamill34/prophet-security-takehome/server/web/internal/db"
	"github.com/jhamill34/prophet-security-takehome/server/web/internal/routes"
	"github.com/jhamill34/prophet-security-takehome/server/web/internal/server"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	queries := database.New(
		db.NewDatabase(context.TODO(), "host=localhost port=5432 user=prophet-th password=prophet-th dbname=prophet-th sslmode=disable"),
	)

	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(swagger, &nethttpmiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: auth.AuthValidator,
		},
		SilenceServersWarning: true,
	})

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(validator)

	// TODO: Custom error functions
	serverRoutes := api.NewStrictHandler(routes.NewServerRoutes(queries), []api.StrictMiddlewareFunc{})

	handler := api.HandlerWithOptions(serverRoutes, api.ChiServerOptions{
		BaseRouter: router,
	})

	s := server.NewHttpServer(3333, handler)

	slog.Info("Starting server")
	s.ListenAndServe()
}
