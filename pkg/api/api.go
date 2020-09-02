package api

import (
	"context"
	"log"
	"net/http"

	"ismacaulay/procrast-api/pkg/db"

	"github.com/go-chi/chi"
)

type Api struct {
	router *chi.Mux
}

func New(db db.Database) *Api {
	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user", "example")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Route("/procrast/v1", func(r chi.Router) {
		r.Route("/lists", func(r chi.Router) {
			r.Get("/", getListsHandler(db))
			r.Post("/", postListHandler)

			r.Route("/{id:[0-9a-zA-Z]+}", func(r chi.Router) {
				r.Get("/", getListHandler)
				r.Patch("/", patchListHandler)
				r.Delete("/", deleteListHandler)

				r.Route("/items", func(r chi.Router) {
					r.Get("/", getItemsHandler)
					r.Post("/", postItemHandler)

					r.Route("/{id:[0-9a-zA-Z]+}", func(r chi.Router) {
						r.Get("/", getItemHandler)
						r.Patch("/", patchItemHandler)
						r.Delete("/", deleteItemHandler)
					})
				})
			})
		})
	})

	return &Api{router: r}
}

func (api *Api) Run() {
	log.Println("Starting...")
	log.Fatal(http.ListenAndServe(":8080", api.router))
	log.Println("Shutting down")
}
