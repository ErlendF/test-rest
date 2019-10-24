package server

import (
	"github.com/gorilla/mux"
)

// NewRouter creates a new router
func newRouter(h *Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/test", h.testHandler)
	r.HandleFunc("/", h.testHandler)

	return r
}
