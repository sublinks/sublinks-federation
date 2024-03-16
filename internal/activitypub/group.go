package activitypub

import (
	"fmt"
	"sublinks/sublinks-federation/internal/model"
	"time"
)

type Group struct {
	Context           *Context  `json:"@context"`
	Id                string    `json:"id"`
	PreferredUsername string    `json:"preferredUsername"`
	Inbox             string    `json:"inbox"`
	Outbox            string    `json:"outbox"`
	Type              string    `json:"type"`
	Summary           string    `json:"summary"`
	Content           string    `json:"content"`
	Image             []Link    `json:"image"`
	Icon              []Link    `json:"icon"`
	Source            Source    `json:"source"`
	Publickey         PublicKey `json:"publicKey"`
	Published         time.Time `json:"published"`
	Endpoints         Endpoints `json:"endpoints"`
}

func NewGroup(
	id string,
	displayName string,
	name string,
	bio string,
	publickey string,
) Group {
	group := Group{}
	group.Context = GetContext()
	group.Id = id
	group.PreferredUsername = name
	group.Inbox = fmt.Sprintf("%s/inbox", id)
	group.Outbox = fmt.Sprintf("%s/outbox", id)
	group.Type = "Group"
	group.Summary = bio
	group.Publickey = PublicKey{
		Keyid:        fmt.Sprintf("%s#main-key", id),
		Owner:        id,
		PublicKeyPem: publickey,
	}
	group.Published = time.Now()
	group.Endpoints.SharedInbox = fmt.Sprintf("%s/inbox", id)
	return group
}

func ConvertActorToGroup(a *model.Actor) Group {
	return NewGroup(
		a.Id,
		a.Username,
		a.Name,
		a.Bio,
		a.PublicKey,
	)
}
