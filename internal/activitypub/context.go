package activitypub

type Context []interface{}

type lemmyContextData struct {
	Lemmy                   string `json:"lemmy"`
	Litepub                 string `json:"litepub"`
	Pt                      string `json:"pt"`
	Sc                      string `json:"sc"`
	ChatMessage             string
	CommentsEnabled         string    `json:"commentsEnabled"`
	Sensitive               string    `json:"sensitive"`
	MatrixUserId            string    `json:"matrixUserId"`
	PostingRestrictedToMods string    `json:"postingRestrictedToMods"`
	RemoveData              string    `json:"removeData"`
	Stickied                string    `json:"stickied"`
	Moderators              moderator `json:"moderators"`
	Expires                 string    `json:"expires"`
	Distinguished           string    `json:"distinguished"`
	Language                string    `json:"language"`
	Identifier              string    `json:"identifier"`
}

type moderator struct {
	Type string `json:"@type"`
	Id   string `json:"@id"`
}

func GetContext() Context {
	return Context{
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1",
		lemmyContextData{
			Lemmy:                   "https://join-lemmy.org/ns#",
			Litepub:                 "http://litepub.social/ns#",
			Pt:                      "https://joinpeertube.org/ns#",
			Sc:                      "http://schema.org/",
			ChatMessage:             "litepub:ChatMessage",
			CommentsEnabled:         "pt:commentsEnabled",
			Sensitive:               "as:sensitive",
			MatrixUserId:            "lemmy:matrixUserId",
			PostingRestrictedToMods: "lemmy:postingRestrictedToMods",
			RemoveData:              "lemmy:removeData",
			Stickied:                "lemmy:stickied",
			Moderators: moderator{
				Type: "@id",
				Id:   "lemmy:moderators",
			},
			Expires:       "as:endTime",
			Distinguished: "lemmy:distinguished",
			Language:      "sc:inLanguage",
			Identifier:    "sc:identifier",
		},
	}
}
