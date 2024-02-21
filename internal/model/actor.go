package model

type Actor struct {
	ActorType  string `json:"actor_type"`
	Id         string `json:"id"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}
