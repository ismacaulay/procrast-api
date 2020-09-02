package api

import (
	"encoding/json"
	"net/http"

	"ismacaulay/procrast-api/pkg/db"
	"ismacaulay/procrast-api/pkg/models"

	"github.com/google/uuid"
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

func postListHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)

		var request models.List
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		err = request.Validate(true)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		id, err := uuid.NewRandom()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		request.Id = &id
		err = db.CreateList(user, request)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusCreated, request)
	}
}

func getListHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func patchListHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func deleteListHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
