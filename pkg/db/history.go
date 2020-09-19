package db

import (
	"ismacaulay/procrast-api/pkg/models"
	"log"

	"github.com/google/uuid"
)

func GetHistorySince(conn Conn, user string, since uint64) ([]models.History, error) {
	sqlStatement := `
		SELECT id, command, state, created
		FROM history
		WHERE user_id = $1 AND created >= $2
		ORDER BY created ASC`

	rows, err := conn.Query(sqlStatement, user, since)
	if err != nil {
		log.Printf("Failed to load history for user %s\nError: %s\n", user, err.Error())
		return []models.History{}, ErrFailedToLoadData
	}
	defer rows.Close()

	history := make([]models.History, 0)
	for rows.Next() {
		var id uuid.UUID
		var command string
		var state []byte
		var created int64
		if err := rows.Scan(&id, &command, &state, &created); err != nil {
			log.Printf("Failed to scan row: %s\n", err.Error())
			return []models.History{}, ErrFailedToScanRow
		}

		item := models.History{
			UUID:    id,
			Command: command,
			State:   state,
			Created: created,
		}
		history = append(history, item)
	}

	return history, nil
}

func GetHistory(conn Conn, user string, historyId uuid.UUID) (models.History, error) {
	sqlStatement := `
		SELECT id, command, state, created
		FROM history
		WHERE user_id = $1 AND id = $2`

	var id uuid.UUID
	var command string
	var state []byte
	var created int64
	err := conn.QueryRow(sqlStatement, user, historyId).Scan(
		&id, &command, &state, &created)
	if err != nil {
		log.Printf("Failed to execute query: %s\n", err.Error())
		return models.History{}, ErrFailedToLoadData
	}

	history := models.History{
		UUID:    id,
		Command: command,
		State:   state,
		Created: created,
	}
	return history, nil
}

func CreateHistory(conn Conn, user string, history models.History) error {
	sqlStatement := `
		INSERT INTO history (id, command, state, created, user_id)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := conn.Exec(sqlStatement,
		history.UUID, history.Command, history.State, history.Created, user)
	if err != nil {
		log.Println("Failed to create history:", err)
		return ErrFailedToInsert
	}

	return nil
}
