package server

import (
	"fmt"
	"net/http"
	"time"
)

// New creates a new http server
func New(port string, apiVer string) *http.Server {
	handler := newHandler()
	router := newRouter(handler, apiVer)

	return &http.Server{
		Addr: fmt.Sprintf(":%s", port),

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Passing mux router as handler
	}
}
