package db

import (
	"log"

	"ismacaulay/procrast-api/pkg/models"

	"github.com/google/uuid"
)

func FindUserByEmail(conn Conn, email string) (models.User, error) {
	sqlStatement := `
		SELECT id, email, passhash, created, modified
		FROM users
		WHERE email = $1`

	var id uuid.UUID
	var storedEmail string
	var passHash []byte
	var created, modified int64
	err := conn.QueryRow(sqlStatement, email).Scan(&id, &storedEmail, &passHash, &created, &modified)
	if err != nil {
		return models.User{}, ErrFailedToLoadData
	}

	user := models.User{
		UUID:     id,
		Email:    storedEmail,
		PassHash: passHash,
		Created:  created,
		Modified: modified,
	}
	return user, nil
}

func FindUserByUUID(conn Conn, id string) (models.User, error) {
	sqlStatement := `
		SELECT id, email, passhash, created, modified
		FROM users
		WHERE id = $1`

	var storedId uuid.UUID
	var storedEmail string
	var passHash []byte
	var created, modified int64
	err := conn.QueryRow(sqlStatement, id).Scan(&storedId, &storedEmail, &passHash, &created, &modified)
	if err != nil {
		return models.User{}, ErrFailedToLoadData
	}

	user := models.User{
		UUID:     storedId,
		Email:    storedEmail,
		PassHash: passHash,
		Created:  created,
		Modified: modified,
	}
	return user, nil
}

func CreateUser(conn Conn, user models.User) error {
	sqlStatement := `
		INSERT INTO users (id, email, passhash, created, modified)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := conn.Exec(sqlStatement, user.UUID, user.Email, user.PassHash, user.Created, user.Modified)
	if err != nil {
		log.Println("Failed to create user:", err)
		return ErrFailedToInsert
	}

	return nil
}
