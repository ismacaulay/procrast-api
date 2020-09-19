package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"ismacaulay/procrast-api/pkg/db"
	"ismacaulay/procrast-api/pkg/models"

	"github.com/google/uuid"
)

var (
	CmdListCreate = "LIST CREATE"
	CmdListUpdate = "LIST UPDATE"
	CmdListDelete = "LIST DELETE"
	CmdItemCreate = "ITEM CREATE"
	CmdItemUpdate = "ITEM UPDATE"
	CmdItemDelete = "ITEM DELETE"
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

		history, err := db.GetHistorySince(conn, user, since)
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

		processed := make([]uuid.UUID, 0)
		for _, history := range *request.History {
			err = db.Transaction(conn, func(tx db.Conn) error {
				if _, err := db.GetHistory(tx, user, history.UUID); err == nil {
					log.Println("Skipping: History already exists", history.UUID)
					processed = append(processed, history.UUID)
					return nil
				}

				if err := db.CreateHistory(tx, user, history); err != nil {
					return err
				}

				switch history.Command {
				case CmdListCreate:
					{
						var state models.List
						if err := json.Unmarshal(history.State, &state); err != nil {
							return err
						}

						if err := db.CreateList(tx, user, state); err != nil {
							return err
						}

						processed = append(processed, state.UUID)
					}
				case CmdListUpdate:
					{
						var state models.List
						if err := json.Unmarshal(history.State, &state); err != nil {
							return err
						}

						list, err := db.RetrieveList(tx, user, state.UUID.String())
						if err != nil {
							return err
						}

						list.Title = state.Title
						list.Description = state.Description
						list.Modified = state.Modified
						if err := db.UpdateList(tx, user, list); err != nil {
							return err
						}

						processed = append(processed, state.UUID)
					}
				case CmdListDelete:
					{
						var state struct {
							UUID uuid.UUID `json:"uuid"`
						}
						if err := json.Unmarshal(history.State, &state); err != nil {
							return err
						}

						list, err := db.RetrieveList(conn, user, state.UUID.String())
						if err != nil {
							return err
						}

						if err := db.DeleteList(tx, user, list); err != nil {
							return err
						}

						processed = append(processed, state.UUID)
					}
				case CmdItemCreate:
					{
						var state models.Item
						if err := json.Unmarshal(history.State, &state); err != nil {
							return err
						}

						if err := db.CreateItem(tx, state); err != nil {
							return err
						}

						processed = append(processed, state.UUID)
					}
				case CmdItemUpdate:
					{
						var state models.Item
						if err := json.Unmarshal(history.State, &state); err != nil {
							return err
						}

						item, err := db.RetrieveItem(tx, user, state.UUID.String())
						if err != nil {
							return err
						}

						item.Title = state.Title
						item.Description = state.Description
						item.State = state.State
						item.Modified = state.Modified
						if err := db.UpdateItem(tx, item); err != nil {
							return err
						}

						processed = append(processed, state.UUID)
					}
				case CmdItemDelete:
					{
						var state struct {
							UUID uuid.UUID `json:"uuid"`
						}
						if err := json.Unmarshal(history.State, &state); err != nil {
							return err
						}

						item, err := db.RetrieveItem(conn, user, state.UUID.String())
						if err != nil {
							return err
						}

						if err := db.DeleteItem(tx, item); err != nil {
							return err
						}

						processed = append(processed, state.UUID)
					}
				}

				return nil
			})

			// if err != nil {
			// 	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			// 	return
			// }
		}

		respondWithJSON(w, http.StatusCreated, struct {
			Processed []uuid.UUID `json:"processed"`
		}{Processed: processed})
	}
}
