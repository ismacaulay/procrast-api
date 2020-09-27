package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

type TokenClaims struct {
	User string `json:"user,omitempty"`
	jwt.StandardClaims
}

func GenerateToken(id string) (string, error) {
	iat := time.Now()
	// TODO: There is probably a better way of doing this
	exp := iat.Add(time.Hour * 24 * 3) // 3 day exp (for now)
	iss, _ := os.Hostname()

	standardClaims := jwt.StandardClaims{ExpiresAt: exp.Unix(), IssuedAt: iat.Unix(), Issuer: iss}
	claims := TokenClaims{id, standardClaims}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(JWT_SECRET)
	if err != nil {
		log.Fatalf("Unable to generate token: %f", err)
		return "", err
	}
	return tokenStr, nil
}

func DecodeToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return JWT_SECRET, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error decoding token: %f", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return token, nil
}

func ExtractClaims(token *jwt.Token) *TokenClaims {
	return token.Claims.(*TokenClaims)
}
