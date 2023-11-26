package db

import (
	"database/sql"
	"embed"
	"log"
	"os"

	_ "github.com/libsql/libsql-client-go/libsql"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func GetDb(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("libsql", dbUrl)
	return db, err
}

func RunMigrations() {
	dbUrl, _ := os.LookupEnv("DB_URL")
	db, err := GetDb(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
