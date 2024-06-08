package actors

import (
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/model"
)

type CommunityService struct {
	db db.Database
}

func NewCommunityService(db db.Database) *CommunityService {
	return &CommunityService{db}
}

func (a CommunityService) GetById(id string) *model.Actor {
	actor := &model.Actor{ActorType: "Group"}
	a.load(actor, id)
	return actor
}

func (a CommunityService) load(actor *model.Actor, id string) {
	_ = a.db.Find(actor, id)
}

func (a CommunityService) Save(actor *model.Actor) bool {
	err := a.db.Save(actor)
	return err == nil
}
