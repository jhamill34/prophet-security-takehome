package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

const (
	PORT = 3334
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	server := NewServer(PORT, mux)
	Run(server)
}

func NewServer(port uint16, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: handler,
	}
}

func Run(server *http.Server) {
	slog.Info("Starting server")
	server.ListenAndServe()
}
