package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates a new router
func newRouter(h *handler) *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(h.notFound)

	r.HandleFunc("/post", h.getPosts).Methods("GET").Name("getPost")
	r.HandleFunc("/post", h.addPost).Methods("POST").Name("addPost")
	r.HandleFunc("/comment", h.addComment).Methods("POST").Name("addComment")

	r.Use(log)

	return r
}
