package worker

import (
	"encoding/json"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/model"
	"sublinks/sublinks-federation/internal/repository"
)

type ActorWorker struct {
	log.Logger
	repository.Repository
}

func (w *ActorWorker) Process(msg []byte) error {
	actor := model.Actor{}
	err := json.Unmarshal(msg, &actor)
	if err != nil {
		w.Logger.Error("Error unmarshalling actor", err)
		return err
	}
	err = w.Repository.Save(&actor)
	if err != nil {
		w.Logger.Error("Error saving actor", err)
		return err
	}
	return nil
}
