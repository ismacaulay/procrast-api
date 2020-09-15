package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ismacaulay/procrast-api/pkg/db"
	"ismacaulay/procrast-api/pkg/models"
)

var (
	CmdListCreate = "LIST CREATE"
	CmdListUpdate = "LIST UPDATE"
	CmdListDelete = "LIST DELETE"
	CmdItemCreate = "ITEM CREATE"
)

func getHistoryHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)
		since_param := r.URL.Query().Get("since")
		since := uint64(0)
		if since_param != "" {
			s, err := strconv.ParseUint(since_param, 10, 64)
			if err != nil {
				respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
				return
			}

			since = s
		}

		history, err := db.GetHistory(conn, user, since)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusOK, struct {
			History []models.History `json:"history"`
		}{History: history})
	}
}

func postHistoryHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(string)

		var request struct {
			History *[]models.History `json:"history,omitempty"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		if request.History == nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		err = db.Transaction(conn, func(tx db.Conn) error {
			for _, history := range *request.History {
				if err := db.CreateHistory(tx, user, history); err != nil {
					return err
				}

				switch history.Command {
				case CmdItemCreate:
					{
						var state models.Item
						if err := json.Unmarshal(history.State, &state); err != nil {
							return err
						}

						if err := db.CreateItem(tx, state); err != nil {
							return err
						}
					}
				}
			}

			return nil
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusCreated, struct {
			Processed int `json:"processed"`
		}{Processed: len(*request.History)})
	}
}
