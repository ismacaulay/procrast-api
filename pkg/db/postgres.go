package db

import "ismacaulay/procrast-api/pkg/models"

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

type PostgresDatabase struct {
}

func NewPostgresDatabase() *PostgresDatabase {
	return &PostgresDatabase{}
}

func (db *PostgresDatabase) RetrieveAllLists(user string) ([]models.List, error) {
	return []models.List{}, nil
}

func (db *PostgresDatabase) RetrieveList(user, id string) (models.List, error) {
	return models.List{}, nil
}

func (db *PostgresDatabase) CreateList(user string, list models.List) error {
	return nil
}

func (db *PostgresDatabase) UpdateList(user string, list models.List) error {
	return nil
}

func (db *PostgresDatabase) DeleteList(user, list_id string) error {
	return nil
}

func (db *PostgresDatabase) RetrieveAllItems(user, list_id string) ([]models.Item, error) {
	return []models.Item{}, nil
}

func (db *PostgresDatabase) RetrieveItem(user, list_id, item_id string) (models.Item, error) {
	return models.Item{}, nil
}

func (db *PostgresDatabase) CreateItem(user, list_id string, item models.Item) error {
	return nil
}

func (db *PostgresDatabase) UpdateItem(user, list_id string, item models.Item) error {
	return nil
}

func (db *PostgresDatabase) DeleteItem(user, list_id, item_id string) error {
	return nil
}
