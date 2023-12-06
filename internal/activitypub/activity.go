package activitypub

type Activity struct {
	Context  Context     `json:"@context"`
	Actor    string      `json:"actor"`
	To       []string    `json:"to"`
	Object   interface{} `json:"object"`
	Cc       []string    `json:"cc"`
	Audience string      `json:"audience"`
	Type     string      `json:"type"`
	Id       string      `json:"id"`
}

func NewActivity(
	id string,
	actType string,
	actor string,
	to []string,
	cc []string,
	audience string,
	obj interface{},
) Activity {
	return Activity{
		Context:  *GetContext(),
		Id:       id,
		Type:     actType,
		Actor:    actor,
		To:       to,
		Cc:       cc,
		Audience: audience,
		Object:   obj,
	}
}
