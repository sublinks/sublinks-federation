package worker

import (
	"encoding/json"
	"errors"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/model"
	"sublinks/sublinks-federation/internal/service"
)

type ActorWorker struct {
	log.Logger
	userService      *service.UserService
	communityService *service.CommunityService
}

func NewActorWorker(logger log.Logger, userService *service.UserService, communityService *service.CommunityService) *ActorWorker {
	return &ActorWorker{
		Logger:           logger,
		userService:      userService,
		communityService: communityService,
	}
}

func (w *ActorWorker) Process(msg []byte) error {
	actor := model.Actor{}
	err := json.Unmarshal(msg, &actor)
	if err != nil {
		w.Logger.Error("Error unmarshalling actor", err)
		return err
	}
	if actor.ActorType == "Group" && !w.communityService.Save(&actor) {
		w.Logger.Error("Error saving actor (community)", nil)
		return errors.New("Error saving actor (community)")
	}

	if actor.ActorType == "Person" && !w.userService.Save(&actor) {
		w.Logger.Error("Error saving actor (user)", nil)
		return errors.New("Error saving actor (user)")
	}
	return nil
}
