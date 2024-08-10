package server

import (
	"fmt"
	"net/http"
)

func NewHttpServer(
	port uint16,
	handler http.Handler,
) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: handler,
	}
}
