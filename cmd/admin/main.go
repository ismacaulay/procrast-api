package main

import (
	"bufio"
	"fmt"
	"ismacaulay/procrast-api/pkg/models"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"ismacaulay/procrast-api/pkg/db"
)

func main() {
	userDbConfig := db.PostgresConfig{
		Host:     os.Getenv("USERDB_HOST"),
		Port:     os.Getenv("USERDB_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB"),
	}
	userDb := db.NewPostgresDatabase(userDbConfig)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)
	if email == "" {
		fmt.Println("No email entered")
		os.Exit(1)
	}

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.Trim(password, "\n")
	if password == "" {
		fmt.Println("No password entered")
		os.Exit(1)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		fmt.Println("Failed to generate password", err)
		os.Exit(1)
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("Failed to generate uuid", err)
		os.Exit(1)
	}

	now := time.Now().UTC().Unix()
	user := models.User{
		UUID:     uuid,
		Email:    email,
		PassHash: passHash,
		Created:  now,
		Modified: now,
	}

	if err := db.CreateUser(userDb.Conn, user); err != nil {
		fmt.Println("Failed to create user", err)
		os.Exit(1)
	}

	fmt.Println("Created user: ", uuid)
}
