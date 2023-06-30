package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ApiHandler() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/ping", pingHandler())
	})
	return r
}
