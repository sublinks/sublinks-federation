package comments

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

func (p CommentService) FindComment(id string) *model.Comment {
	comment := &model.Comment{}
	err := p.db.Find(comment, id)
	if err != nil {
		return comment
	}
	return nil
}

func (p *CommentService) Save(comment *model.Comment) bool {
	err := p.db.Save(comment)
	return err == nil
}
