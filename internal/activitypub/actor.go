package activitypub

type PublicKey struct {
	Keyid        string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}

type Endpoints struct {
	SharedInbox string `json:"sharedInbox"`
}
