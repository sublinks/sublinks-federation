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
	// bootstrap logger
	logger := log.NewLogger("main")

	// Load connection string from .env file
	err := godotenv.Load()
	if err != nil {
		logger.Warn(fmt.Sprintf("failed to load env, %v", err))
	}

	conn, err := db.Connect()
	if err != nil {
		logger.Fatal("failed connecting to db", err)
	}
	defer conn.Close()
	db.RunMigrations(conn)

	s := http.NewServer(logger)
	s.RunServer()

	os.Exit(0)
}
