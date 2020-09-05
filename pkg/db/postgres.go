package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"ismacaulay/procrast-api/pkg/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Database interface {
	RetrieveAllLists(user string) ([]models.List, error)
	RetrieveList(user, id string) (models.List, error)
	CreateList(user string, list models.List) error
	UpdateList(user string, list models.List) error
	DeleteList(user, id string) error

	RetrieveAllItems(user, list_id string) ([]models.Item, error)
	RetrieveItem(user, list_id, item_id string) (models.Item, error)
	CreateItem(user, list_id string, item models.Item) error
	UpdateItem(user, list_id string, item models.Item) error
	DeleteItem(user, list_id, item_id string) error
}

var (
	ErrFailedToLoadData   = errors.New("Failed to load data")
	ErrFailedToUpdateData = errors.New("Failed to update data")
	ErrFailedToDeleteData = errors.New("Failed to delete data")
	ErrFailedToScanRow    = errors.New("Failed to scan row")
	ErrFailedToInsert     = errors.New("Failed to insert into database")
	ErrPasswordMismatch   = errors.New("Failed to validate password")
)

var (
	host     = os.Getenv("POSTGRES_HOST")
	port     = os.Getenv("POSTGRES_PORT")
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

type PostgresDatabase struct {
	conn *sql.DB
}

func NewPostgresDatabase() *PostgresDatabase {
	log.Println("Initializing postgres")
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	conn, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalf("Unable to open db: %s", err)
	}

	if err = conn.Ping(); err != nil {
		log.Fatalf("Unable to connect to db: %s", err)
	}

	return &PostgresDatabase{conn}
}

func (db *PostgresDatabase) RetrieveAllLists(user string) ([]models.List, error) {
	sqlStatement := `
		SELECT id, title, description FROM lists
		WHERE user_id = $1 ORDER BY date_created DESC`

	rows, err := db.conn.Query(sqlStatement, user)
	if err != nil {
		log.Printf("Failed to load messages for user %s\nError: %s\n", user, err.Error())
		return []models.List{}, ErrFailedToLoadData
	}
	defer rows.Close()

	lists := make([]models.List, 0)
	for rows.Next() {
		var id uuid.UUID
		var title, description string
		if err := rows.Scan(&id, &title, &description); err != nil {
			log.Printf("Failed to scan row: %s\n", err.Error())
			return []models.List{}, ErrFailedToScanRow
		}

		list := models.List{
			Id:          &id,
			Title:       &title,
			Description: &description,
		}
		lists = append(lists, list)
	}

	return lists, nil
}

func (db *PostgresDatabase) RetrieveList(user, id string) (models.List, error) {
	sqlStatement := `SELECT id, title, description FROM lists WHERE user_id = $1 AND id = $2`

	var listId uuid.UUID
	var title, description string
	err := db.conn.QueryRow(sqlStatement, user, id).Scan(&listId, &title, &description)
	if err != nil {
		log.Printf("Failed to execute query: %s\n", err.Error())
		return models.List{}, ErrFailedToLoadData
	}

	list := models.List{
		Id:          &listId,
		Title:       &title,
		Description: &description,
	}
	return list, nil
}

func (db *PostgresDatabase) CreateList(user string, list models.List) error {
	now := time.Now().UTC().Unix()

	sqlStatement := `
		INSERT INTO lists (id, date_created, date_modified, title, description, user_id)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.conn.Exec(sqlStatement, list.Id, now, now, list.Title, list.Description, user)
	if err != nil {
		log.Println("Failed to create list:", err)
		return ErrFailedToInsert
	}

	return nil
}

func (db *PostgresDatabase) UpdateList(user string, list models.List) error {
	now := time.Now().UTC().Unix()

	sqlStatement := `
		UPDATE lists
		SET date_modified = $3, title = $4, description = $5
		WHERE user_id = $1 AND id = $2`

	_, err := db.conn.Exec(sqlStatement, user, list.Id, now, list.Title, list.Description)
	if err != nil {
		log.Println("Failed to update list:", err)
		return ErrFailedToUpdateData
	}

	return nil
}

func (db *PostgresDatabase) DeleteList(user, list_id string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		return ErrFailedToDeleteData
	}

	sqlStatement := `
		DELETE FROM items i
		USING lists l
		WHERE l.user_id = $1 AND l.id = i.list_id AND l.id = $2`

	_, err = tx.Exec(sqlStatement, user, list_id)
	if err != nil {
		log.Println("Failed to delete items from list:", err)
		tx.Rollback()
		return ErrFailedToDeleteData
	}

	sqlStatement = `
		DELETE FROM lists
		WHERE user_id = $1 AND id = $2`

	_, err = tx.Exec(sqlStatement, user, list_id)
	if err != nil {
		log.Println("Failed to delete list:", err)
		tx.Rollback()
		return ErrFailedToDeleteData
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %s\n", err.Error())
		return ErrFailedToDeleteData
	}
	return nil
}

func (db *PostgresDatabase) RetrieveAllItems(user, list_id string) ([]models.Item, error) {
	sqlStatement := `
		SELECT i.id, i.title, i.description
		FROM items i
		INNER JOIN lists l ON (i.list_id = l.id)
		WHERE l.user_id = $1 AND i.list_id = $2
		ORDER BY i.date_created DESC`

	rows, err := db.conn.Query(sqlStatement, user, list_id)
	if err != nil {
		log.Printf("Failed to load items for user %s\nError: %s\n", user, err.Error())
		return []models.Item{}, ErrFailedToLoadData
	}
	defer rows.Close()

	items := make([]models.Item, 0)
	for rows.Next() {
		var id uuid.UUID
		var title, description string
		if err := rows.Scan(&id, &title, &description); err != nil {
			log.Printf("Failed to scan row: %s\n", err.Error())
			return []models.Item{}, ErrFailedToScanRow
		}

		item := models.Item{
			Id:          &id,
			Title:       &title,
			Description: &description,
		}
		items = append(items, item)
	}

	return items, nil
}

func (db *PostgresDatabase) RetrieveItem(user, list_id, item_id string) (models.Item, error) {
	sqlStatement := `
		SELECT i.id, i.title, i.description
		FROM items i
		INNER JOIN lists l ON (i.list_id = l.id)
		WHERE l.user_id = $1 AND i.list_id = $2 AND i.id = $3
		ORDER BY i.date_created DESC`

	var listId uuid.UUID
	var title, description string
	err := db.conn.QueryRow(sqlStatement, user, list_id, item_id).Scan(&listId, &title, &description)
	if err != nil {
		log.Printf("Failed to execute query: %s\n", err.Error())
		return models.Item{}, ErrFailedToLoadData
	}

	item := models.Item{
		Id:          &listId,
		Title:       &title,
		Description: &description,
	}
	return item, nil
}

func (db *PostgresDatabase) CreateItem(user, list_id string, item models.Item) error {
	now := time.Now().UTC().Unix()

	sqlStatement := `
		INSERT INTO items (id, date_created, date_modified, title, description, list_id)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.conn.Exec(sqlStatement, item.Id, now, now, item.Title, item.Description, list_id)
	if err != nil {
		log.Println("Failed to create item:", err)
		return ErrFailedToInsert
	}

	return nil
}

func (db *PostgresDatabase) UpdateItem(user, list_id string, item models.Item) error {
	now := time.Now().UTC().Unix()

	sqlStatement := `
		UPDATE items
		SET date_modified = $4, title = $5, description = $6
		FROM lists
		WHERE lists.user_id = $1
			AND items.list_id = $2
			AND items.id = $3`

	_, err := db.conn.Exec(sqlStatement,
		user, list_id, item.Id, now, item.Title, item.Description,
	)
	if err != nil {
		log.Println("Failed to update item:", err)
		return ErrFailedToUpdateData
	}

	return nil
}

func (db *PostgresDatabase) DeleteItem(user, list_id, item_id string) error {
	sqlStatement := `
		DELETE FROM items
		USING lists
		WHERE lists.user_id = $1
			AND items.list_id = $2
			AND items.id = $3`

	_, err := db.conn.Exec(sqlStatement, user, list_id, item_id)
	if err != nil {
		log.Println("Failed to update item:", err)
		return ErrFailedToUpdateData
	}
	return nil
}
