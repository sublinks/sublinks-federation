package db

import (
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	"sublinks/sublinks-federation/internal/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed "migrations"
var migrations embed.FS

func RunMigrations(db *sql.DB) {
	log.Debug("Running migrations...")
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal("Error getting MySQL driver", err)
	}
	source, _ := iofs.New(migrations, "migrations")
	m, err := migrate.NewWithInstance("iofs", source, "mysql", driver)
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}
	if err := m.Up(); err != nil && fmt.Sprintf("%s", err) != "no change" {
		log.Fatal("Error running migrations", err)
	}
	log.Debug("Done!")
}
