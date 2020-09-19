package db

import (
	"log"

	"ismacaulay/procrast-api/pkg/models"

	"github.com/google/uuid"
)

const selectAllItemsStatement = `
	SELECT i.id, i.title, i.description, i.state, i.created, i.modified, l.id
	FROM items i
	INNER JOIN lists l ON (i.list_id = l.id)
	WHERE l.user_id = $1 AND i.list_id = $2
	ORDER BY i.created ASC`

const selectItemStatement = `
	SELECT i.id, i.title, i.description, i.state, i.created, i.modified, l.id
	FROM items i
	INNER JOIN lists l ON (i.list_id = l.id)
	WHERE l.user_id = $1 AND i.id = $2
	ORDER BY i.created ASC`

const insertItemStatement = `
	INSERT INTO items (id, created, modified, title, description, state, list_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

const deleteItemStatement = `
	DELETE FROM items
	USING lists
	WHERE items.id = $1`

func RetrieveAllItems(conn Conn, user, list_id string) ([]models.Item, error) {
	rows, err := conn.Query(selectAllItemsStatement, user, list_id)
	if err != nil {
		log.Printf("Failed to load items for user %s\nError: %s\n", user, err.Error())
		return []models.Item{}, ErrFailedToLoadData
	}
	defer rows.Close()

	items := make([]models.Item, 0)
	for rows.Next() {
		var itemId, listId uuid.UUID
		var title, description string
		var state uint8
		var created, modified int64
		if err := rows.Scan(&itemId, &title, &description, &state, &created, &modified, &listId); err != nil {
			log.Printf("Failed to scan row: %s\n", err.Error())
			return []models.Item{}, ErrFailedToScanRow
		}

		item := models.Item{
			UUID:        itemId,
			Title:       title,
			Description: description,
			State:       state,
			Created:     created,
			Modified:    modified,
			ListUUID:    listId,
		}
		items = append(items, item)
	}

	return items, nil
}

func RetrieveItem(conn Conn, user, id string) (models.Item, error) {
	var itemId, listId uuid.UUID
	var title, description string
	var state uint8
	var created, modified int64
	err := conn.QueryRow(selectItemStatement, user, id).Scan(
		&itemId, &title, &description, &state, &created, &modified, &listId)
	if err != nil {
		log.Printf("Failed to execute query: %s\n", err.Error())
		return models.Item{}, ErrFailedToLoadData
	}

	item := models.Item{
		UUID:        itemId,
		Title:       title,
		Description: description,
		State:       state,
		Created:     created,
		Modified:    modified,
		ListUUID:    listId,
	}
	return item, nil
}

func CreateItem(conn Conn, item models.Item) error {
	_, err := conn.Exec(insertItemStatement,
		item.UUID, item.Created, item.Modified,
		item.Title, item.Description, item.State, item.ListUUID)
	if err != nil {
		log.Println("Failed to create item:", err)
		return ErrFailedToInsert
	}

	return nil
}

func UpdateItem(conn Conn, item models.Item) error {
	sqlStatement := `
		UPDATE items
		SET modified = $2, title = $3, description = $4, state = $5
		FROM lists
		WHERE items.id = $1`

	_, err := conn.Exec(sqlStatement,
		item.UUID, item.Modified, item.Title, item.Description, item.State,
	)
	if err != nil {
		log.Println("Failed to update item:", err)
		return ErrFailedToUpdateData
	}

	return nil
}

func DeleteItem(conn Conn, item models.Item) error {
	_, err := conn.Exec(deleteItemStatement, item.UUID)
	if err != nil {
		log.Println("Failed to update item:", err)
		return ErrFailedToUpdateData
	}
	return nil
}
