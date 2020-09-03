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

	// TODO: This is a temporary to add a user id to the request.
	//       The user ID should come from the token
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
			r.Post("/", postListHandler(db))

			r.Route("/{listId:[0-9a-zA-Z]+}", func(r chi.Router) {
				r.Get("/", getListHandler(db))
				r.Patch("/", patchListHandler(db))
				r.Delete("/", deleteListHandler(db))

				r.Route("/items", func(r chi.Router) {
					r.Get("/", getItemsHandler(db))
					r.Post("/", postItemHandler(db))

					r.Route("/{itemId:[0-9a-zA-Z]+}", func(r chi.Router) {
						r.Get("/", getItemHandler(db))
						r.Patch("/", patchItemHandler(db))
						r.Delete("/", deleteItemHandler(db))
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
