package db

import (
	"ismacaulay/procrast-api/pkg/models"
	"log"

	"github.com/google/uuid"
)

func RetrieveAllLists(conn Conn, user string) ([]models.List, error) {
	sqlStatement := `
		SELECT id, title, description, created, modified FROM lists
		WHERE user_id = $1 ORDER BY created DESC`

	rows, err := conn.Query(sqlStatement, user)
	if err != nil {
		log.Printf("Failed to load messages for user %s\nError: %s\n", user, err.Error())
		return []models.List{}, ErrFailedToLoadData
	}
	defer rows.Close()

	lists := make([]models.List, 0)
	for rows.Next() {
		var id uuid.UUID
		var title, description string
		var created, modified int64
		if err := rows.Scan(&id, &title, &description, &created, &modified); err != nil {
			log.Printf("Failed to scan row: %s\n", err.Error())
			return []models.List{}, ErrFailedToScanRow
		}

		list := models.List{
			UUID:        id,
			Title:       title,
			Description: description,
			Created:     created,
			Modified:    modified,
		}
		lists = append(lists, list)
	}

	return lists, nil
}

func RetrieveList(conn Conn, user, id string) (models.List, error) {
	sqlStatement := `SELECT id, title, description, created, modified FROM lists WHERE user_id = $1 AND id = $2`

	var listId uuid.UUID
	var title, description string
	var created, modified int64
	err := conn.QueryRow(sqlStatement, user, id).Scan(&listId, &title, &description, &created, &modified)
	if err != nil {
		log.Printf("Failed to execute query: %s\n", err.Error())
		return models.List{}, ErrFailedToLoadData
	}

	list := models.List{
		UUID:        listId,
		Title:       title,
		Description: description,
		Created:     created,
		Modified:    modified,
	}
	return list, nil
}

func CreateList(conn Conn, user string, list models.List) error {
	sqlStatement := `
		INSERT INTO lists (id, created, modified, title, description, user_id)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := conn.Exec(sqlStatement, list.UUID, list.Created, list.Modified, list.Title, list.Description, user)
	if err != nil {
		log.Println("Failed to create list:", err)
		return ErrFailedToInsert
	}

	return nil
}

func UpdateList(conn Conn, user string, list models.List) error {
	sqlStatement := `
		UPDATE lists
		SET modified = $3, title = $4, description = $5
		WHERE user_id = $1 AND id = $2`

	_, err := conn.Exec(sqlStatement, user, list.UUID, list.Modified, list.Title, list.Description)
	if err != nil {
		log.Println("Failed to update list:", err)
		return ErrFailedToUpdateData
	}

	return nil
}

func DeleteList(conn Conn, user string, list models.List) error {
	sqlStatement := `
		DELETE FROM items i
		USING lists l
		WHERE l.user_id = $1 AND l.id = i.list_id AND l.id = $2`

	_, err := conn.Exec(sqlStatement, user, list.UUID)
	if err != nil {
		log.Println("Failed to delete items from list:", err)
		return ErrFailedToDeleteData
	}

	sqlStatement = `
		DELETE FROM lists
		WHERE user_id = $1 AND id = $2`

	_, err = conn.Exec(sqlStatement, user, list.UUID)
	if err != nil {
		log.Println("Failed to delete list:", err)
		return ErrFailedToDeleteData
	}

	return nil
}
