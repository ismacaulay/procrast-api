package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func validateUUIDParameterMiddleware(param string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, param)
			_, err := uuid.Parse(id)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
