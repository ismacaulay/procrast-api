package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ismacaulay/procrast-api/pkg/db"
	"ismacaulay/procrast-api/pkg/models"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func getItemsHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		listId := chi.URLParam(r, "listId")

		if _, err := db.RetrieveList(conn, user, listId); err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		items, err := db.RetrieveAllItems(conn, user, listId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusOK, struct {
			Items []models.Item `json:"items"`
		}{Items: items})
	}
}

func postItemHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		listId := chi.URLParam(r, "listId")
		now := time.Now().UTC().Unix()

		list, err := db.RetrieveList(conn, user, listId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		var request struct {
			Title       *string `json:"title,omitempty"`
			Description string  `json:"description"`
		}
		if json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		if request.Title == nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		id, err := uuid.NewRandom()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		item := models.Item{
			UUID:        id,
			Title:       *request.Title,
			Description: request.Description,
			Created:     now,
			Modified:    now,
			ListUUID:    list.UUID,
		}

		err = db.Transaction(conn, func(conn db.Conn) error {
			if err := db.CreateItem(conn, item); err != nil {
				return err
			}

			state, err := json.Marshal(item)
			if err != nil {
				return err
			}

			id, err := uuid.NewRandom()
			if err != nil {
				return err
			}
			history := models.History{
				UUID:    id,
				Command: CmdItemCreate,
				State:   state,
				Created: now,
			}

			if err := db.CreateHistory(conn, user, history); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			log.Println("Failed to execute transaction:", err)
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusCreated, item)
	}
}

func getItemHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		itemId := chi.URLParam(r, "itemId")

		item, err := db.RetrieveItem(conn, user, itemId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		respondWithJSON(w, http.StatusOK, item)
	}
}

func patchItemHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		itemId := chi.URLParam(r, "itemId")

		var request struct {
			Title       *string `json:"title,omitempty"`
			Description *string `json:"description,omitempty"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		item, err := db.RetrieveItem(conn, user, itemId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		if request.Title != nil {
			item.Title = *request.Title
		}

		if request.Description != nil {
			item.Description = *request.Description
		}

		if db.UpdateItem(conn, user, item) != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusOK, item)
	}
}

func deleteItemHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		itemId := chi.URLParam(r, "itemId")

		item, err := db.RetrieveItem(conn, user, itemId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		if db.DeleteItem(conn, user, item) != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusNoContent, nil)
	}
}
