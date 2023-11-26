package lemmy

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
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

func GetLemmyClient(ctx context.Context) *Client {
	user, _ := os.LookupEnv("LEMMY_USER")
	pass, _ := os.LookupEnv("LEMMY_PASSWORD")
	return NewClient("https://demo.sublinks.org", user, pass)
}

func (c *Client) GetPost(ctx context.Context, id string) (*PostResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v3/post?id=%s", c.BaseURL, id), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	res := PostResponse{}
	if _, err := c.sendRequest(req, &res); err != nil {
		return nil, errors.New(err.Message)
	}
	return &res, nil
}

func (c *Client) GetUser(ctx context.Context, id string) (*UserResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v3/user?username=%s", c.BaseURL, id), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	res := UserResponse{}
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

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
