package server

import (
	"github.com/gorilla/mux"
)

// NewRouter creates a new router
func newRouter(h *Handler, apiVer string) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/test", h.testHandler)

	return r
}
