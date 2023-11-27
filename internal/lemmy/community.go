package lemmy

import "time"

type CommunityCounts struct {
	Id                  int    `json:"id"`
	CommunityId         int    `json:"community_id"`
	Subscribers         int    `json:"subscribers"`
	Posts               int    `json:"posts"`
	Comments            int    `json:"comments"`
	Published           string `json:"published"`
	UsersActiveDay      int    `json:"users_active_day"`
	UsersActiveWeek     int    `json:"users_active_week"`
	UsersActiveMonth    int    `json:"users_active_month"`
	UsersActiveHalfYear int    `json:"users_active_half_year"`
	Hot_rank            int    `json:"hot_rank"`
}

type CommunityView struct {
	Community  Community       `json:"community"`
	Subscribed string          `json:"subscribed"`
	Blocked    bool            `json:"blocked"`
	Counts     CommunityCounts `json:"counts"`
}

type Community struct {
	Id                      int       `json:"id"`
	Name                    string    `json:"name"`
	Title                   string    `json:"title"`
	Description             string    `json:"description"`
	Removed                 bool      `json:"removed"`
	Published               time.Time `json:"published"`
	Updated                 time.Time `json:"updated"`
	Deleted                 bool      `json:"deleted"`
	Nsfw                    bool      `json:"nsfw"`
	ActorId                 string    `json:"actor_id"`
	Local                   bool      `json:"local"`
	Icon                    string    `json:"icon"`
	Banner                  string    `json:"banner"`
	FollowersUrl            string    `json:"followers_url"`
	InboxUrl                string    `json:"inbox_url"`
	Hidden                  bool      `json:"hidden"`
	PostingRestrictedToMods bool      `json:"posting_restricted_to_mods"`
	InstanceId              int       `json:"instance_id"`
}
