package db

import (
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/model"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (d *FederationDB) RunMigrations() {
	logger := log.NewLogger("db migrations")
	logger.Debug("Running migrations...")
	err := d.DB.AutoMigrate(&model.Actor{}, &model.Post{})
	if err != nil {
		logger.Fatal("Failed to run migrations", err)
	}
	logger.Debug("Done!")
}
