package repository

import (
	"sublinks/sublinks-federation/internal/db"
)

type Repository interface {
	Save(interface{}) error
	Find(*interface{}, ...interface{}) error
}

type RepositoryImpl struct {
	db db.Database
}

func NewRepository(db db.Database) Repository {
	return &RepositoryImpl{db: db}
}

func (repository *RepositoryImpl) Save(a interface{}) error {
	return repository.db.Save(a)
}

func (repository *RepositoryImpl) Find(a *interface{}, params ...interface{}) error {
	return repository.db.Find(a, params...)
}
