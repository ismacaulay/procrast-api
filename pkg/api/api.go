package api

import (
	"log"
	"net/http"

	"ismacaulay/procrast-api/pkg/auth"
	"ismacaulay/procrast-api/pkg/db"

	"github.com/go-chi/chi"
)

type Api struct {
	router *chi.Mux
}

func New(db, userDb db.DB) *Api {
	r := chi.NewRouter()

	r.Get("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/auth/v1", func(r chi.Router) {
		r.Post("/login", postLoginHandler(userDb))
	})

	r.Route("/procrast/v1", func(r chi.Router) {
		r.Use(auth.TokenSecurity)
		r.Use(auth.UserValidation(userDb))

		r.Route("/lists", func(r chi.Router) {
			r.Get("/", getListsHandler(db))
			r.Post("/", postListHandler(db))

			r.Route("/{listId}", func(r chi.Router) {
				r.Use(validateUUIDParameterMiddleware("listId"))

				r.Get("/", getListHandler(db))
				r.Patch("/", patchListHandler(db))
				r.Delete("/", deleteListHandler(db))

				r.Get("/items", getItemsHandler(db))
				r.Post("/items", postItemHandler(db))
			})
		})

		r.Route("/items", func(r chi.Router) {
			r.Route("/{itemId}", func(r chi.Router) {
				r.Use(validateUUIDParameterMiddleware("itemId"))

				r.Get("/", getItemHandler(db))
				r.Patch("/", patchItemHandler(db))
				r.Delete("/", deleteItemHandler(db))
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
