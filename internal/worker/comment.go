package worker

import (
	"encoding/json"
	"os"
	"strings"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/model"
	"sublinks/sublinks-federation/internal/repository"
)

type CommentWorker struct {
	log.Logger
	repository.Repository
}

func (w *CommentWorker) Process(msg []byte) error {
	comment := model.Comment{}
	err := json.Unmarshal(msg, &comment)
	hostname := os.Getenv("HOSTNAME")
	comment.UrlStub = strings.Replace(hostname, comment.Id, "", 1)
	if err != nil {
		w.Logger.Error("Error unmarshalling post: %s", err)
		return err
	}
	err = w.Repository.Save(comment)
	if err != nil {
		w.Logger.Error("Error saving post: %s", err)
		return err
	}
	return nil
}
