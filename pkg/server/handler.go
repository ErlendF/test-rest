package server

import (
	"fmt"
	"net/http"
)

//Handler test
type Handler struct{}

//newHandler returns handler
func newHandler() *Handler {
	return &Handler{}
}

func (h *Handler) testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Success! Change test 21.10 17.10 \n")
}
