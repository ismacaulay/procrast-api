package api

import (
	"net/http"
	"strconv"

	"ismacaulay/procrast-api/pkg/db"
	"ismacaulay/procrast-api/pkg/models"
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
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}
