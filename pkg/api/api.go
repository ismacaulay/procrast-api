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

func New(db db.DB) *Api {
	r := chi.NewRouter()

	// TODO: This is a temporary to add a user id to the request.
	//       The user ID should come from the token
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userId := "1aaa6956-ac0a-4500-96a5-91803dcf8894"
			ctx = context.WithValue(ctx, "user", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Route("/procrast/v1", func(r chi.Router) {
		r.Route("/lists", func(r chi.Router) {
			r.Get("/", getListsHandler(db))
			r.Post("/", postListHandler(db))

			r.Route("/{listId}", func(r chi.Router) {
				r.Use(validateUUIDParameterMiddleware("listId"))

				r.Get("/", getListHandler(db))
				r.Patch("/", patchListHandler(db))
				r.Delete("/", deleteListHandler(db))

				r.Route("/items", func(r chi.Router) {
					r.Get("/", getItemsHandler(db))
					r.Post("/", postItemHandler(db))

					r.Route("/{itemId}", func(r chi.Router) {
						r.Use(validateUUIDParameterMiddleware("itemId"))

						r.Get("/", getItemHandler(db))
						r.Patch("/", patchItemHandler(db))
						r.Delete("/", deleteItemHandler(db))
					})
				})
			})
		})

		r.Route("/history", func(r chi.Router) {
			r.Get("/", getHistoryHandler(db))
			r.Post("/", postHistoryHandler(db))
		})
	})

	return &Api{router: r}
}

func (api *Api) Run() {
	log.Println("Starting...")
	log.Fatal(http.ListenAndServe(":8080", api.router))
	log.Println("Shutting down")
}
