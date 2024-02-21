package model

type Post struct {
	Id      string `json:"id"`
	Action  string `json:"action"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
}
