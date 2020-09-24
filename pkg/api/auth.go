package api

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"ismacaulay/procrast-api/pkg/auth"
	"ismacaulay/procrast-api/pkg/db"
)

func postLoginHandler(conn db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Email    *string `json:"email,omitempty"`
			Password *string `json:"password,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		if request.Email == nil || request.Password == nil {
			respondWithError(w, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
			return
		}

		user, err := db.FindUserByEmail(conn, *request.Email)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(*request.Password)); err != nil {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		token, err := auth.GenerateToken(user.UUID.String())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusOK, struct {
			Token string `json:"token"`
		}{Token: token})
	}
}
