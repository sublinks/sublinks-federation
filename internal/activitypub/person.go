package activitypub

import (
	"fmt"
	"sublinks/sublinks-federation/internal/model"
	"time"
)

type Person struct {
	Context           *Context  `json:"@context"`
	Id                string    `json:"id"`
	PreferredUsername string    `json:"preferredUsername"`
	Inbox             string    `json:"inbox"`
	Outbox            string    `json:"outbox"`
	Type              string    `json:"type"`
	Summary           string    `json:"summary"`
	MatrixUserId      string    `json:"matrixUserId"`
	Image             []Link    `json:"image"`
	Icon              []Link    `json:"icon"`
	Source            Source    `json:"source"`
	Publickey         PublicKey `json:"publicKey"`
	Published         time.Time `json:"published"`
	Endpoints         Endpoints `json:"endpoints"`
}

func NewPerson(
	id string,
	name string,
	matrixUserId string,
	bio string,
	publickey string,
) Person {
	person := Person{}
	person.Context = GetContext()
	person.Id = id
	person.PreferredUsername = name
	person.Inbox = fmt.Sprintf("%s/inbox", id)
	person.Outbox = fmt.Sprintf("%s/outbox", id)
	person.Type = "Person"
	person.Summary = bio
	person.MatrixUserId = matrixUserId
	person.Publickey = PublicKey{
		Keyid:        fmt.Sprintf("%s#main-key", id),
		Owner:        id,
		PublicKeyPem: publickey,
	}
	person.Published = time.Now()
	person.Endpoints.SharedInbox = fmt.Sprintf("%s/inbox", id)
	return person
}

func ConvertActorToPerson(u *model.Actor) Person {
	return NewPerson(
		u.Id,
		u.Username,
		u.MatrixUserId,
		u.Bio,
		u.PublicKey,
	)
}
