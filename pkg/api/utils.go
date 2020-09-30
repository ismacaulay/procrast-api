package api

import (
	"encoding/json"
	"ismacaulay/procrast-api/pkg/db"
	"ismacaulay/procrast-api/pkg/models"
	"net/http"

	"github.com/google/uuid"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, status int, msg string) {
	respondWithJSON(w, status, map[string]string{"message": msg})
}

func createHistoryForState(conn db.Conn, cmd, user string, now int64, state interface{}) error {
	encoded, err := json.Marshal(state)
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	history := models.History{
		UUID:      id,
		Command:   cmd,
		State:     encoded,
		Timestamp: now,
		Created:   now,
	}
	if err := db.CreateHistory(conn, user, history); err != nil {
		return err
	}

	return nil
}
