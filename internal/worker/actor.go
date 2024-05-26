package worker

import (
	"encoding/json"
	"errors"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/model"
	"sublinks/sublinks-federation/internal/service/actors"
)

type ActorWorker struct {
	log.Logger
	Service actors.ActorService
}

func (w *ActorWorker) Process(msg []byte) error {
	actor := model.Actor{}
	err := json.Unmarshal(msg, &actor)
	if err != nil {
		w.Logger.Error("Error unmarshalling actor", err)
		return err
	}
	if !w.Service.Save(&actor) {
		w.Logger.Error("Error saving actor", nil)
		return errors.New("Error saving actor")
	}
	return nil
}
