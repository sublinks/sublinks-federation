package activitypub

import (
	"fmt"
	"os"
	"participating-online/sublinks-federation/internal/lemmy"
	"time"
)

var host, _ = os.LookupEnv("HOSTNAME")
var domain, _ = os.LookupEnv("CSB_BASE_PREVIEW_HOST")
var Hostname string = fmt.Sprintf("%s-8080.%s", host, domain)

type PublicKey struct {
	Keyid        string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}

type Endpoints struct {
	SharedInbox string `json:"sharedInbox"`
}

type User struct {
	Context           Context   `json:"@context"`
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

func NewUser(name string, matrixUserId string, bio string, publickey string) User {
	user := User{}
	user.Context = GetContext()
	user.Id = fmt.Sprintf("https://%s/users/%s", Hostname, name)
	user.PreferredUsername = name
	user.Inbox = fmt.Sprintf("https://%s/users/%s/inbox", Hostname, name)
	user.Outbox = fmt.Sprintf("https://%s/users/%s/outbox", Hostname, name)
	user.Type = "Person"
	user.Summary = bio
	user.MatrixUserId = matrixUserId
	owner := fmt.Sprintf("https://%s/users/%s", Hostname, name)
	user.Publickey = PublicKey{
		Keyid:        fmt.Sprintf("https://%s/users/%s#main-key", Hostname, name),
		Owner:        owner,
		PublicKeyPem: publickey,
	}
	user.Published = time.Now()
	user.Endpoints.SharedInbox = fmt.Sprintf("https://%s/inbox", Hostname)
	return user
}

func ConvertUserToApub(u *lemmy.UserResponse) User {
	return NewUser(
		u.PersonView.Person.Name,
		u.PersonView.Person.MatrixUserId,
		u.PersonView.Person.Bio,
		"", //TODO: Public key goes here
	)
}
