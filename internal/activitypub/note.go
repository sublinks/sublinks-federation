package activitypub

import (
	"fmt"
	"sublinks/sublinks-federation/internal/model"
	"time"
)

type Note struct {
	Context       *Context  `json:"@context,omitempty"`
	Id            string    `json:"id"`
	Type          string    `json:"type"`
	AttributedTo  string    `json:"attributedTo"`
	To            []string  `json:"to"`
	Cc            []string  `json:"cc"`
	Audience      string    `json:"audience"`
	InReplyTo     string    `json:"inReplyTo"`
	Content       string    `json:"content"`
	MediaType     string    `json:"mediaType"`
	Source        Source    `json:"source,omitempty"`
	Tag           []Tag     `json:"tag,omitempty"`
	Distinguished bool      `json:"distinguished,omitempty"`
	Language      Language  `json:"language,omitempty"`
	Published     time.Time `json:"published"`
	Updated       time.Time `json:"updated"`
}

type Tag struct {
	Href string `json:"href"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func NewNote(commentUrl string, fromUser string, postUrl string, commentBody string, published time.Time) *Note {
	return &Note{
		Id:           commentUrl,
		Type:         "Note",
		AttributedTo: fromUser,
		To:           []string{"https://www.w3.org/ns/activitystreams#Public"},
		Cc:           []string{fromUser, commentUrl},
		Audience:     commentUrl,
		InReplyTo:    postUrl,
		Content:      commentBody,
		MediaType:    "text/html",
		Source: Source{
			Content:   fmt.Sprintf("This is a comment on %s post", postUrl),
			MediaType: "text/markdown",
		},
		Language: Language{
			Identifier: "en",
			Name:       "English",
		},
		Distinguished: false,
		Published:     published,
	}
}

func ConvertCommentToNote(c *model.Comment) *Note {
	return NewNote(c.UrlStub, c.Author, c.Post, c.Content, c.Published)
}
