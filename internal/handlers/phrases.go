package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func HandlerWithoutId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "Get:HandlerWithoutId")
	case http.MethodPost:
		fmt.Fprintln(w, "Post:HandlerWithoutId")
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func HandlerWithId(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	args := strings.Split(path, "/")
	id, err := strconv.Atoi(args[len(args)-1])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "Get:HandlerWithId:%d", id)
	case http.MethodDelete:
		fmt.Fprintf(w, "Delete:HandlerWithId:%d", id)
	case http.MethodPatch:
		fmt.Fprintf(w, "Patch:HandlerWithId:%d", id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
