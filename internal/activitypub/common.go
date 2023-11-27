package activitypub

type Source struct {
	Content   string `json:"content"`
	MediaType string `json:"mediaType"`
}

type Link struct {
	Type string `json:"type"` //"Link" | "Image"
	Href string `json:"href"` //"https://enterprise.lemmy.ml/pictrs/image/eOtYb9iEiB.png"
}
