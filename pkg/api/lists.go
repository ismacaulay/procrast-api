package api

import (
	"ismacaulay/procrast-api/pkg/db"
	"ismacaulay/procrast-api/pkg/models"
	"net/http"
)

func getListsHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)

		lists, err := db.RetrieveAllLists(user)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		respondWithJSON(w, http.StatusOK, struct {
			Lists []models.List `json:"lists"`
		}{Lists: lists})
	}
}

func postListHandler(w http.ResponseWriter, r *http.Request) {
}

func getListHandler(w http.ResponseWriter, r *http.Request) {
}

func patchListHandler(w http.ResponseWriter, r *http.Request) {
}

func deleteListHandler(w http.ResponseWriter, r *http.Request) {
}
