package activitypub

import (
	"fmt"
	"time"
)

type Source struct {
	Content   string `json:"content"`
	MediaType string `json:"mediaType"`
}

type Link struct {
	Type string `json:"type"` //"Link" | "Image"
	Href string `json:"href"` //"https://enterprise.lemmy.ml/pictrs/image/eOtYb9iEiB.png"
}

type Language struct {
	Identifier string `json:"identifier"` // "fr",
	Name       string `json:"name"`       // "Français"
}

type Post struct {
	Context         Context   `json:"@context"`
	Id              string    `json:"id"`
	Type            string    `json:"type"`
	AttributedTo    string    `json:"attributedTo"`
	To              []string  `json:"to"`
	Cc              []string  `json:"cc"`
	Audience        string    `json:"audience"`
	Name            string    `json:"name"`
	Content         string    `json:"content"`
	MediaType       string    `json:"mediaType"`
	Source          Source    `json:"source"`
	Attachment      []Link    `json:"attachment"`
	Image           []Link    `json:"image"`
	Sensitive       bool      `json:"sensitive"`
	CommentsEnabled bool      `json:"commentsEnabled"`
	Language        Language  `json:"language"`
	Published       time.Time `json:"published"`
}

func NewPost(postUrl string, fromUser string, communityUrl string, postTitle string, postBody string, nsfw bool, published time.Time) Post {
	post := Post{
		Context:      GetContext(),
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
		//Attachment
		//Image
		Sensitive:       nsfw,
		CommentsEnabled: true,
		Language: Language{
			Identifier: "en",
			Name:       "English",
		},
		Published: published,
	}
	return post
}
