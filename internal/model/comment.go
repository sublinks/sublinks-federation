package model

import "time"

type Comment struct {
	Id        string    `json:"id" gorm:"primary_key"`
	UrlStub   string    `json:"url_stub"`
	Post      string    `json:"post_id"`
	Author    string    `json:"author"`
	Nsfw      bool      `json:"nsfw"`
	Published time.Time `json:"published"`
	Content   string    `json:"content"`
}
