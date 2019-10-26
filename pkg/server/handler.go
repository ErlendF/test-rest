package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test/pkg/database"
	"test/pkg/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//handler test
type handler struct {
	database.Database
}

//newHandler returns handler
func newHandler(db *database.Database) *handler {
	return &handler{*db}
}

func (h *handler) testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Success! Change test 19.10 \n")
}

func (h *handler) addPost(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("route", mux.CurrentRoute(r).GetName()).Debugf("Request recieved")
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName()}).Warn("Could not parse request body")
		http.Error(w, "Bad request: bad request body", http.StatusBadRequest)
		return
	}
	err = h.AddPost(post.Content)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName()}).Warn("Could not add post")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OK")
}

func (h *handler) addComment(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("route", mux.CurrentRoute(r).GetName()).Debugf("Request recieved")
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName()}).Warn("Could not parse request body")
		http.Error(w, "Bad request: bad request body", http.StatusBadRequest)
		return
	}
	err = h.AddComment(&comment)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName()}).Warn("Could not add comment")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OK")
}

func (h *handler) getPosts(w http.ResponseWriter, r *http.Request) {
	logrus.WithField("route", mux.CurrentRoute(r).GetName()).Debugf("Request recieved")
	resp, err := h.GetPosts()
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName()}).Warn("Could not get posts")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logrus.Debugf("resp: %+v", resp)
	setHeader(w)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"route": mux.CurrentRoute(r).GetName(), "response": resp}).Warn("Could not encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
