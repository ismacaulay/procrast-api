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

func getListsHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)

		lists, err := db.RetrieveAllLists(conn, user)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		respondWithJSON(w, http.StatusOK, struct {
			Lists []models.List `json:"lists"`
		}{Lists: lists})
	}
}

func postListHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		now := time.Now().UTC().Unix()

		var request struct {
			Title       *string `json:"title,omitempty"`
			Description string  `json:"description"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
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

		list := models.List{
			UUID:        id,
			Title:       *request.Title,
			Description: request.Description,
			Created:     now,
			Modified:    now,
		}

		err = db.Transaction(conn, func(tx db.Conn) error {
			if err := db.CreateList(tx, user, list); err != nil {
				return err
			}

			state, err := json.Marshal(list)
			if err != nil {
				return err
			}

			id, err := uuid.NewRandom()
			if err != nil {
				return err
			}

			history := models.History{
				UUID:    id,
				Command: CmdListCreate,
				State:   state,
				Created: now,
			}
			if err := db.CreateHistory(tx, user, history); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Println("Failed to execute transaction:", err)
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusCreated, list)
	}
}

func getListHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		listId := chi.URLParam(r, "listId")

		list, err := db.RetrieveList(conn, user, listId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		respondWithJSON(w, http.StatusOK, list)
	}
}

func patchListHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		listId := chi.URLParam(r, "listId")

		var request struct {
			Title       *string `json:"title,omitempty"`
			Description *string `json:"description,omitempty"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		list, err := db.RetrieveList(conn, user, listId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		update := false
		if request.Title != nil {
			list.Title = *request.Title
			update = true
		}

		if request.Description != nil {
			list.Description = *request.Description
			update = true
		}

		if update {
			now := time.Now().UTC().Unix()
			list.Modified = now

			err = db.Transaction(conn, func(tx db.Conn) error {
				if err = db.UpdateList(tx, user, list); err != nil {
					return err
				}

				state, err := json.Marshal(list)
				if err != nil {
					return err
				}

				id, err := uuid.NewRandom()
				if err != nil {
					return err
				}

				history := models.History{
					UUID:    id,
					Command: CmdListUpdate,
					State:   state,
					Created: now,
				}
				if err := db.CreateHistory(tx, user, history); err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}

		}

		respondWithJSON(w, http.StatusOK, list)
	}
}

func deleteListHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		listId := chi.URLParam(r, "listId")
		list, err := db.RetrieveList(conn, user, listId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		err = db.Transaction(conn, func(tx db.Conn) error {
			if err := db.DeleteList(tx, user, list); err != nil {
				return err
			}
			state, err := json.Marshal(struct {
				UUID uuid.UUID `json:"uuid"`
			}{UUID: list.UUID})
			if err != nil {
				return err
			}

			id, err := uuid.NewRandom()
			if err != nil {
				return err
			}

			now := time.Now().UTC().Unix()
			history := models.History{
				UUID:    id,
				Command: CmdListDelete,
				State:   state,
				Created: now,
			}
			if err := db.CreateHistory(tx, user, history); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusNoContent, nil)
	}
}
