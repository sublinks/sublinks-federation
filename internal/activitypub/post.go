package activitypub

import (
	"fmt"
	"sublinks/sublinks-federation/internal/model"
	"time"
)

type Page struct {
	Context         *Context  `json:"@context,omitempty"`
	Id              string    `json:"id"`
	Type            string    `json:"type"`
	AttributedTo    string    `json:"attributedTo"`
	To              []string  `json:"to"`
	Cc              []string  `json:"cc,omitempty"`
	Audience        string    `json:"audience"`
	Name            string    `json:"name"`
	Content         string    `json:"content"`
	MediaType       string    `json:"mediaType"`
	Source          Source    `json:"source"`
	Attachment      []Link    `json:"attachment,omitempty"`
	Image           []Link    `json:"image,omitempty"`
	Sensitive       bool      `json:"sensitive"`
	CommentsEnabled bool      `json:"commentsEnabled"`
	Language        Language  `json:"language"`
	Published       time.Time `json:"published"`
}

func NewPage(postUrl string, fromUser string, communityUrl string, postTitle string, postBody string, nsfw bool, published time.Time) *Page {
	return &Page{
		Id:           postUrl,
		Type:         "Page",
		AttributedTo: fromUser,
		To:           []string{"https://www.w3.org/ns/activitystreams#Public", communityUrl},
		Cc:           []string{},
		Audience:     communityUrl,
		Name:         postTitle,
		Content:      postBody,
		MediaType:    "text/html",
		Source: Source{
			Content:   fmt.Sprintf("This is a post in the %s community", communityUrl),
			MediaType: "text/markdown",
		},
		// Attachment
		// Image
		Sensitive:       nsfw,
		CommentsEnabled: true,
		Language: Language{
			Identifier: "en",
			Name:       "English",
		},
		Published: published,
	}
}

func ConvertPostToPage(p *model.Post) *Page {
	return NewPage(
		p.Id,
		p.Author,
		p.Community,
		p.Title,
		p.Content,
		p.Nsfw,
		p.Published,
	)
}
