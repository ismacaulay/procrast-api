package api

import (
	"encoding/json"
	"net/http"

	"ismacaulay/procrast-api/pkg/db"
	"ismacaulay/procrast-api/pkg/models"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func getItemsHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		listId := chi.URLParam(r, "listId")

		if _, err := db.RetrieveList(user, listId); err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		items, err := db.RetrieveAllItems(user, listId)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		respondWithJSON(w, http.StatusOK, struct {
			Items []models.Item `json:"items"`
		}{Items: items})
	}
}

func postItemHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		listId := chi.URLParam(r, "listId")

		if _, err := db.RetrieveList(user, listId); err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		var request models.Item
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		if request.Validate(true) != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		id, err := uuid.NewRandom()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		request.Id = &id
		err = db.CreateItem(user, listId, request)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusCreated, request)
	}
}

func getItemHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		listId := chi.URLParam(r, "listId")
		itemId := chi.URLParam(r, "itemId")

		if _, err := db.RetrieveList(user, listId); err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		item, err := db.RetrieveItem(user, listId, itemId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		respondWithJSON(w, http.StatusOK, item)
	}
}

func patchItemHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		listId := chi.URLParam(r, "listId")
		itemId := chi.URLParam(r, "itemId")

		if _, err := db.RetrieveList(user, listId); err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		var request models.Item
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		if request.Validate(false) != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		item, err := db.RetrieveItem(user, listId, itemId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		if request.Title != nil {
			item.Title = request.Title
		}

		if request.Description != nil {
			item.Description = request.Description
		}

		if db.UpdateItem(user, listId, item) != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusOK, item)
	}
}

func deleteItemHandler(db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		listId := chi.URLParam(r, "listId")
		itemId := chi.URLParam(r, "itemId")

		if _, err := db.RetrieveList(user, listId); err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		if db.DeleteItem(user, listId, itemId) != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusNoContent, nil)
	}
}
