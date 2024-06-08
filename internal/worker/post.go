package worker

import (
	"encoding/json"
	"os"
	"strings"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/model"
	"sublinks/sublinks-federation/internal/service"
)

type PostWorker struct {
	log.Logger
	service *service.PostService
}

func NewPostWorker(logger log.Logger, service *service.PostService) *PostWorker {
	return &PostWorker{
		Logger:  logger,
		service: service,
	}
}

func (w *PostWorker) Process(msg []byte) error {
	post := model.Post{}
	err := json.Unmarshal(msg, &post)
	hostname := os.Getenv("HOSTNAME")
	post.UrlStub = strings.Replace(hostname, post.Id, "", 1)
	if err != nil {
		w.Logger.Error("Error unmarshalling post: %s", err)
		return err
	}
	res := w.service.Save(&post)
	if !res {
		w.Logger.Error("Error saving post", nil)
		return err
	}
	return nil
}
