package server

import (
	"fmt"
	"net/http"
	"test/pkg/database"
	"time"
)

// New creates a new http server
func New(db *database.Database, port int) *http.Server {
	handler := newHandler(db)
	router := newRouter(handler)

	return &http.Server{
		Addr: fmt.Sprintf(":%d", port),

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Passing mux router as handler
	}
}
