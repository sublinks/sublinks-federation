package worker

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/model"
	"sublinks/sublinks-federation/internal/service"
)

type CommentWorker struct {
	log.Logger
	service *service.CommentService
}

func NewCommentWorker(logger log.Logger, service *service.CommentService) *CommentWorker {
	return &CommentWorker{
		Logger:  logger,
		service: service,
	}
}

func (w *CommentWorker) Process(msg []byte) error {
	comment := model.Comment{}
	err := json.Unmarshal(msg, &comment)
	hostname := os.Getenv("HOSTNAME")
	comment.UrlStub = strings.Replace(hostname, comment.Id, "", 1)
	if err != nil {
		w.Logger.Error("Error unmarshalling comment: %s", err)
		return err
	}
	if !w.service.Save(&comment) {
		w.Logger.Error("Error saving comment", nil)
		return errors.New("Error saving comment")
	}
	return nil
}
