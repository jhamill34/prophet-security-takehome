package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Resource interface {
	Path() string
	Handler() http.Handler
}

func RootRouter(resources ...Resource) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	for _, resource := range resources {
		r.Mount(resource.Path(), resource.Handler())
	}
	return r
}

func NewHttpServer(
	port uint16,
	handler http.Handler,
) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: handler,
	}
}
