package actors

import (
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/model"
)

type ActorService struct {
	db db.Database
}

func NewActorService(db db.Database) *ActorService {
	return &ActorService{db}
}

func (a ActorService) FindCommunity(id string) *model.Actor {
	actor := model.Actor{ActorType: "Group"}
	return a.Load(&actor, id)
}

func (a ActorService) FindUser(id string) *model.Actor {
	actor := model.Actor{ActorType: "Person"}
	return a.Load(&actor, id)
}

func (a ActorService) Load(actor *model.Actor, id string) *model.Actor {
	err := a.db.Find(actor, id)
	if err != nil {
		return actor
	}
	return nil
}

func (a *ActorService) Save(actor *model.Actor) bool {
	err := a.db.Save(actor)
	return err == nil
}
