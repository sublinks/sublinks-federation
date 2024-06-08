package service

import (
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/model"
)

type CommentService struct {
	db db.Database
}

func NewCommentService(db db.Database) *CommentService {
	return &CommentService{db}
}

func (p CommentService) GetById(id string) interface{} {
	comment := &model.Comment{}
	err := p.db.Find(comment, id)
	if err != nil {
		return comment
	}
	return nil
}

func (p CommentService) Save(comment interface{}) bool {
	err := p.db.Save(comment)
	return err == nil
}
