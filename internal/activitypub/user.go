package activitypub

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
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
	Context           string    `json:"@context"`
	Id                string    `json:"id"`
	PreferredUsername string    `json:"preferredUsername"`
	Inbox             string    `json:"inbox"`
	Outbox            string    `json:"outbox"`
	Type              string    `json:"type"`
	Name              string    `json:"name"`
	Publickey         PublicKey `json:"publicKey"`
	privatekey        string
	Published         time.Time `json:"published"`
	Updated           time.Time `json:"updated"`
	Endpoints         Endpoints `json:"endpoints"`
}

func GetPrivateKeyString(privatekey *rsa.PrivateKey) string {
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privatekey)
	if err != nil {
		log.Fatalf("error when dumping privatekey: %s \n", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}))
}

func NewExistingUser(name string, privatekey string, publickey string) User {
	user := User{}
	user.Context = "https://www.w3.org/ns/activitystreams"
	user.PreferredUsername = name
	user.Id = fmt.Sprintf("https://%s/users/%s", Hostname, name)
	user.Inbox = fmt.Sprintf("https://%s/users/%s/inbox", Hostname, name)
	user.Outbox = fmt.Sprintf("https://%s/users/%s/outbox", Hostname, name)
	user.Type = "Person"
	user.Name = name
	user.privatekey = privatekey
	pubKeyId := fmt.Sprintf("https://%s/users/%s#main-key", Hostname, name)
	owner := fmt.Sprintf("https://%s/users/%s", Hostname, name)
	user.Publickey = PublicKey{Keyid: pubKeyId, PublicKeyPem: publickey, Owner: owner}
	user.Published = time.Now()
	user.Updated = time.Now()
	user.Endpoints.SharedInbox = fmt.Sprintf("https://%s/inbox", Hostname)
	return user
}

func NewUser(name string, privatekey *rsa.PrivateKey, publickey *rsa.PublicKey) User {
	user := User{}
	user.Context = "https://www.w3.org/ns/activitystreams"
	user.Id = fmt.Sprintf("https://%s/users/%s", Hostname, name)
	user.PreferredUsername = name
	user.Inbox = fmt.Sprintf("https://%s/users/%s/inbox", Hostname, name)
	user.Outbox = fmt.Sprintf("https://%s/users/%s/outbox", Hostname, name)
	user.Type = "Person"
	user.Name = name
	user.privatekey = GetPrivateKeyString(privatekey)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		fmt.Printf("error when dumping publickey: %s \n", err)
		return User{}
	}
	owner := fmt.Sprintf("https://%s/users/%s", Hostname, name)
	user.Publickey = PublicKey{
		Keyid: fmt.Sprintf("https://%s/users/%s#main-key", Hostname, name),
		Owner: owner,
		PublicKeyPem: string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		})),
	}
	user.Published = time.Now()
	user.Updated = time.Now()
	user.Endpoints.SharedInbox = fmt.Sprintf("https://%s/inbox", Hostname)
	return user
}
