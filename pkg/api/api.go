package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Api struct {
	router *chi.Mux
}

func New() *Api {
	r := chi.NewRouter()

	r.Route("/procrast/v1", func(r chi.Router) {
		r.Route("/lists", func(r chi.Router) {
			r.Get("/", getListsHandler)
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
