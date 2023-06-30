package api

import (
	"fmt"
	"net/http"
)

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		//nolint:gocritic,revive // linter doesn't like fmt.Fprint(w, "OK\n") but it's fine
		fmt.Fprint(w, "OK\n")
	}
}
