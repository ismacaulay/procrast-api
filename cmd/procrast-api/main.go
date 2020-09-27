package main

import (
	"os"

	"ismacaulay/procrast-api/pkg/api"
	"ismacaulay/procrast-api/pkg/db"
)

func main() {
	dbConfig := db.PostgresConfig{
		Host:     os.Getenv("PROCRASTDB_HOST"),
		Port:     os.Getenv("PROCRASTDB_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB"),
	}
	dataDb := db.NewPostgresDatabase(dbConfig)

	userDbConfig := db.PostgresConfig{
		Host:     os.Getenv("USERDB_HOST"),
		Port:     os.Getenv("USERDB_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB"),
	}
	userDb := db.NewPostgresDatabase(userDbConfig)

	api := api.New(dataDb.Conn, userDb.Conn)
	api.Run()
}
