package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

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
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return tx.Commit()
}
