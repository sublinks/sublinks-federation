package main

import (
	"os"
	"participating-online/sublinks-federation/internal/db"
	"participating-online/sublinks-federation/internal/http"
)

func main() {
	db.RunMigrations()
	http.RunServer()
	os.Exit(0)
}
