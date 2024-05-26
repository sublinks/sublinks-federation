package worker

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/model"
	"sublinks/sublinks-federation/internal/service/comments"
)

type CommentWorker struct {
	log.Logger
	Service comments.CommentService
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
	if !w.Service.Save(&comment) {
		w.Logger.Error("Error saving comment", nil)
		return errors.New("Error saving comment")
	}
	return nil
}
