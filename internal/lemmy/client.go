package lemmy

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	User       string
	Password   string
	HTTPClient *http.Client
}

func NewClient(url string, user string, password string) *Client {
	return &Client{
		BaseURL:  url,
		User:     user,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) GetPost(ctx context.Context, id string) (*Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v3/post?id=%s", c.BaseURL, id), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	res := Response{}
	if _, err := c.sendRequest(req, &res); err != nil {
		return nil, errors.New(err.Message)
	}
	return &res, nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) (*successResponse, *errorResponse) {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.User, c.Password)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, &errorResponse{Code: 500, Message: err.Error()}
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return nil, &errRes
		}
		return nil, &errorResponse{Code: res.StatusCode, Message: fmt.Sprintf("unknown error, status code: %d", res.StatusCode)}
	}

	fullResponse := successResponse{
		Data: v,
	}
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, &errorResponse{Code: 500, Message: err.Error()}
	}
	return &fullResponse, nil
}

type Response struct {
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

type Moderator struct {
	Community Community `json:"community"`
	Moderator Person    `json:"moderator"`
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

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
