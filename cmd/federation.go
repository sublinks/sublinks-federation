package main

import (
	"fmt"
	"os"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/http"
	"sublinks/sublinks-federation/internal/log"

	"github.com/joho/godotenv"
)

func main() {
	// Load connection string from .env file
	err := godotenv.Load()
	if err != nil {
		log.Warn(fmt.Sprintf("failed to load env, %v", err))
	}
	conn, err := db.Connect()
	if err != nil {
		log.Fatal("failed connecting to db", err)
	}
	defer conn.Close()
	db.RunMigrations(conn)
	http.RunServer()
	os.Exit(0)
}
