package db

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Database interface {
	Connect() error
	RunMigrations()
	Close() error
}

type PostgresDB struct {
	*sql.DB
}

func NewDatabase() Database {
	return &PostgresDB{}
}

func (d *PostgresDB) Connect() error {
	database, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		return err
	}
	database.SetConnMaxLifetime(time.Minute * 3)
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(10)
	d.DB = database
	return nil
}
