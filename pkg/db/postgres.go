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

type Conn interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type DB interface {
	Conn

	Begin() (*sql.Tx, error)
}

var (
	ErrFailedToLoadData         = errors.New("Failed to load data")
	ErrFailedToUpdateData       = errors.New("Failed to update data")
	ErrFailedToDeleteData       = errors.New("Failed to delete data")
	ErrFailedToScanRow          = errors.New("Failed to scan row")
	ErrFailedToInsert           = errors.New("Failed to insert into database")
	ErrPasswordMismatch         = errors.New("Failed to validate password")
	ErrFailedToStartTransaction = errors.New("Failed to start transaction")
)

var (
	host     = os.Getenv("POSTGRES_HOST")
	port     = os.Getenv("POSTGRES_PORT")
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

type PostgresDatabase struct {
	Conn *sql.DB
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

func Transaction(conn DB, f func(Conn) error) error {
	tx, err := conn.Begin()
	if err != nil {
		return ErrFailedToStartTransaction
	}

	if err = f(tx); err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

func RetrieveAllLists(conn Conn, user string) ([]models.List, error) {
	sqlStatement := `
		SELECT id, title, description FROM lists
		WHERE user_id = $1 ORDER BY date_created DESC`

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

func RetrieveList(conn Conn, user, id string) (models.List, error) {
	sqlStatement := `SELECT id, title, description FROM lists WHERE user_id = $1 AND id = $2`

	var listId uuid.UUID
	var title, description string
	err := conn.QueryRow(sqlStatement, user, id).Scan(&listId, &title, &description)
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

func CreateList(conn Conn, user string, list models.List) error {
	now := time.Now().UTC().Unix()

	sqlStatement := `
		INSERT INTO lists (id, date_created, date_modified, title, description, user_id)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := conn.Exec(sqlStatement, list.Id, now, now, list.Title, list.Description, user)
	if err != nil {
		log.Println("Failed to create list:", err)
		return ErrFailedToInsert
	}

	return nil
}

func UpdateList(conn Conn, user string, list models.List) error {
	now := time.Now().UTC().Unix()

	sqlStatement := `
		UPDATE lists
		SET date_modified = $3, title = $4, description = $5
		WHERE user_id = $1 AND id = $2`

	_, err := conn.Exec(sqlStatement, user, list.Id, now, list.Title, list.Description)
	if err != nil {
		log.Println("Failed to update list:", err)
		return ErrFailedToUpdateData
	}

	return nil
}

func DeleteList(conn Conn, user, list_id string) error {
	sqlStatement := `
		DELETE FROM items i
		USING lists l
		WHERE l.user_id = $1 AND l.id = i.list_id AND l.id = $2`

	_, err := conn.Exec(sqlStatement, user, list_id)
	if err != nil {
		log.Println("Failed to delete items from list:", err)
		return ErrFailedToDeleteData
	}

	sqlStatement = `
		DELETE FROM lists
		WHERE user_id = $1 AND id = $2`

	_, err = conn.Exec(sqlStatement, user, list_id)
	if err != nil {
		log.Println("Failed to delete list:", err)
		return ErrFailedToDeleteData
	}

	return nil
}

func RetrieveAllItems(conn Conn, user, list_id string) ([]models.Item, error) {
	sqlStatement := `
		SELECT i.id, i.title, i.description
		FROM items i
		INNER JOIN lists l ON (i.list_id = l.id)
		WHERE l.user_id = $1 AND i.list_id = $2
		ORDER BY i.date_created DESC`

	rows, err := conn.Query(sqlStatement, user, list_id)
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

func RetrieveItem(conn Conn, user, list_id, item_id string) (models.Item, error) {
	sqlStatement := `
		SELECT i.id, i.title, i.description
		FROM items i
		INNER JOIN lists l ON (i.list_id = l.id)
		WHERE l.user_id = $1 AND i.list_id = $2 AND i.id = $3
		ORDER BY i.date_created DESC`

	var listId uuid.UUID
	var title, description string
	err := conn.QueryRow(sqlStatement, user, list_id, item_id).Scan(&listId, &title, &description)
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

func CreateItem(conn Conn, user, list_id string, item models.Item) error {
	sqlStatement := `
		INSERT INTO items (id, date_created, date_modified, title, description, list_id)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := conn.Exec(sqlStatement,
		item.Id, item.Created, item.Modified,
		item.Title, item.Description, list_id)
	if err != nil {
		log.Println("Failed to create item:", err)
		return ErrFailedToInsert
	}

	return nil
}

func UpdateItem(conn Conn, user, list_id string, item models.Item) error {
	now := time.Now().UTC().Unix()

	sqlStatement := `
		UPDATE items
		SET date_modified = $4, title = $5, description = $6
		FROM lists
		WHERE lists.user_id = $1
			AND items.list_id = $2
			AND items.id = $3`

	_, err := conn.Exec(sqlStatement,
		user, list_id, item.Id, now, item.Title, item.Description,
	)
	if err != nil {
		log.Println("Failed to update item:", err)
		return ErrFailedToUpdateData
	}

	return nil
}

func DeleteItem(conn Conn, user, list_id, item_id string) error {
	sqlStatement := `
		DELETE FROM items
		USING lists
		WHERE lists.user_id = $1
			AND items.list_id = $2
			AND items.id = $3`

	_, err := conn.Exec(sqlStatement, user, list_id, item_id)
	if err != nil {
		log.Println("Failed to update item:", err)
		return ErrFailedToUpdateData
	}
	return nil
}

func GetHistory(conn Conn, user string, since uint64) ([]models.History, error) {
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
			Id:      id,
			Command: command,
			State:   state,
			Created: created,
		}
		history = append(history, item)
	}

	return history, nil
}

func CreateHistory(conn Conn, user string, history models.History) error {
	sqlStatement := `
		INSERT INTO history (id, command, state, created, user_id)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := conn.Exec(sqlStatement,
		history.Id, history.Command, history.State, history.Created, user)
	if err != nil {
		log.Println("Failed to create history:", err)
		return ErrFailedToInsert
	}

	return nil
}
