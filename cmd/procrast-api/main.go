package main

import (
	"ismacaulay/procrast-api/pkg/api"
	"ismacaulay/procrast-api/pkg/db"
)

func main() {
	db := db.NewPostgresDatabase()

	api := api.New(db)
	api.Run()
}
