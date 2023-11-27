package lemmy

import (
	"time"
)

type PostResponse struct {
	PostView      PostView      `json:"post_view"`
	CommunityView CommunityView `json:"community_view"`
	Moderators    []Moderator   `json:"moderators"`
	CrossPosts    []Post        `json:"cross_posts"`
}

type PostCounts struct {
	Id                     int       `json:"id"`
	PostId                 int       `json:"post_id"`
	Comments               int       `json:"comments"`
	Score                  int       `json:"score"`
	Upvotes                int       `json:"upvotes"`
	Downvotes              int       `json:"downvotes"`
	Published              time.Time `json:"published"`
	NewestCommentTimeNecro time.Time `json:"newest_comment_time_necro"`
	NewestCommentTime      time.Time `json:"newest_comment_time"`
	FeaturedCommunity      bool      `json:"featured_community"`
	FeaturedLocal          bool      `json:"featured_local"`
	HotRank                int       `json:"hot_rank"`
	HotRankActive          int       `json:"hot_rank_active"`
}

type PostView struct {
	Post                       Post       `json:"post"`
	Creator                    Person     `json:"creator"`
	Community                  Community  `json:"community"`
	CreatorBannedFromCommunity bool       `json:"creator_banned_from_community"`
	Counts                     PostCounts `json:"counts"`
	Subscribed                 string     `json:"subscribed"`
	Saved                      bool       `json:"saved"`
	Read                       bool       `json:"read"`
	CreatorBlocked             bool       `json:"creator_blocked"`
	MyVote                     int        `json:"my_vote"`
	UnreadComments             int        `json:"unread_comments"`
}

type Post struct {
	Id                int       `json:"id"`
	Name              string    `json:"name"`
	Url               string    `json:"url"`
	Body              string    `json:"body"`
	CreatorId         int       `json:"creator_id"`
	CommunityId       int       `json:"community_id"`
	Removed           bool      `json:"removed"`
	Locked            bool      `json:"locked"`
	Published         time.Time `json:"published"`
	Updated           time.Time `json:"update"`
	Deleted           bool      `json:"deleted"`
	Nsfw              bool      `json:"nsfw"`
	EmbedTitle        string    `json:"embed_title"`
	EmbedDescription  string    `json:"embed_description"`
	ThumbnailUrl      string    `json:"thumbnail_url"`
	ApId              string    `json:"ap_id"`
	Local             bool      `json:"local"`
	EmbedVideoUrl     string    `json:"embed_video_url"`
	LanguageId        int       `json:"language_id"`
	FeaturedCommunity bool      `json:"featured_community"`
	FeaturedLocal     bool      `json:"featured_local"`
}
