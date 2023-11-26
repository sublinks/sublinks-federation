package lemmy

import "time"

type UserResponse struct {
	PersonView    PersonView    `json:"person_view"`
	CommunityView CommunityView `json:"community_view"`
	Moderators    []Moderator   `json:"moderators"`
	CrossPosts    []Post        `json:"cross_posts"`
}

type PersonView struct {
	Person       Person       `json:"person"`
	PersonCounts PersonCounts `json:"counts"`
	Comments []Comment `json:"comments"`
	Posts []Post `json:"posts"`
	Moderates []Community `json:"moderates"`
}

type PersonCounts struct {
	Id           int `json:"id"`
	PersonId     int `json:"person_id"`
	PostCount    int `json:"post_count"`
	PostScore    int `json:"post_score"`
	CommentCount int `json:"comment_count"`
	CommentScore int `json:"comment_score"`
}

type Moderator struct {
	Community Community `json:"community"`
	Moderator Person    `json:"moderator"`
}

type Person struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	DisplayName    string    `json:"display_name"`
	Avatar         string    `json:"avatar"`
	Banned         bool      `json:"banned"`
	Published      time.Time `json:"published"`
	Updated        time.Time `json:"updated"`
	ActorId        string    `json:"actor_id"`
	Bio            string    `json:"bio"`
	Local          bool      `json:"local"`
	Banner         string    `json:"banner"`
	Deleted        bool      `json:"deleted"`
	InboxUrl       string    `json:"inbox_url"`
	SharedInboxUrl string    `json:"shared_inbox_url"`
	MatrixUserId   string    `json:"matrix_user_id"`
	Admin          bool      `json:"admin"`
	BotAccount     bool      `json:"bot_account"`
	BanExpires     string    `json:"ban_expires"`
	InstanceId     int       `json:"instance_id"`
}
