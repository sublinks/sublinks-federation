package model

import "time"

type Post struct {
	Id        string    `json:"id" gorm:"primary_key"`
	UrlStub   string    `json:"url_stub"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Community string    `json:"community"`
	Nsfw      bool      `json:"nsfw"`
	Published time.Time `json:"published"`
	Content   string    `json:"content"`
}
