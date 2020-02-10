package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"test/pkg/database"
	"test/pkg/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var errBadBody = errors.New("bad request body")

// handler contains a database
type handler struct {
	database.Database
}

// newHandler returns handler
func newHandler(db *database.Database) *handler {
	return &handler{*db}
}

func (h *handler) addPost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		logRespond(w, r, fmt.Errorf("%w: %v", errBadBody, err))
		return
	}
	err = h.AddPost(post.Content)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) addComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		logRespond(w, r, fmt.Errorf("%w: %v", errBadBody, err))
		return
	}
	err = h.AddComment(&comment)
	if err != nil {
		logRespond(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) getPosts(w http.ResponseWriter, r *http.Request) {
	resp, err := h.GetPosts()
	if err != nil {
		logRespond(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logRespond(w, r, err)
		return
	}
}

// notFound handles all requests which don't hit any of the routes defined in the router
func (h *handler) notFound(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("request", r.RequestURI).Debug("Not found handler")
	http.Error(w, "Test-Rest: "+http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// logRespond handles errors. It logs the error and returns an appropriate errormessage and status code based on the error.
func logRespond(w http.ResponseWriter, r *http.Request, err error) {
	logrus.WithField("route", mux.CurrentRoute(r).GetName()).Warn(err)

	switch {
	case errors.Is(err, models.ErrDBInsert):
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	case errors.Is(err, models.ErrNotFound):
		http.Error(w, "No posts yet", http.StatusNotFound)
	case errors.Is(err, errBadBody):
		http.Error(w, "Bad request: bad request body", http.StatusBadRequest)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
