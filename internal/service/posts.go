package service

import (
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/model"
)

type PostService struct {
	db db.Database
}

func NewPostService(db db.Database) *PostService {
	return &PostService{db}
}

func (p PostService) GetById(id string) *model.Post {
	post := &model.Post{}
	err := p.db.Find(post, id)
	if err != nil {
		return post
	}
	return nil
}

func (p PostService) Save(post *model.Post) bool {
	err := p.db.Save(post)
	return err == nil
}
