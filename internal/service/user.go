package service

import (
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/model"
)

type UserService struct {
	db db.Database
}

func NewUserService(db db.Database) *UserService {
	return &UserService{db}
}

func (a UserService) GetById(id string) *model.Actor {
	actor := model.Actor{ActorType: "Person"}
	return a.Load(&actor, id)
}

func (a UserService) Load(actor *model.Actor, id string) *model.Actor {
	err := a.db.Find(actor, id)
	if err != nil {
		return actor
	}
	return nil
}

func (a UserService) Save(actor *model.Actor) bool {
	err := a.db.Save(actor)
	return err == nil
}
