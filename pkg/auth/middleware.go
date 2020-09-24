package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"ismacaulay/procrast-api/pkg/db"
)

func TokenSecurity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := ExtractBearerToken(r)
		if tokenStr == "" {
			log.Println("Missing token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		token, err := DecodeToken(tokenStr)
		if err != nil {
			log.Println("Invalid token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		claims := ExtractClaims(token)
		userId := claims.User
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserValidation(conn db.Conn) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uuid := r.Context().Value("user").(string)

			_, err := db.FindUserByUUID(conn, uuid)
			if err != nil {
				log.Println("Invalid user id from token")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ExtractBearerToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToLower(bearer[0:6]) == "bearer" {
		return bearer[7:]
	}
	return ""
}
