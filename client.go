package main

import (
	"fmt"
	"net/http"

	twitter "github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

type client struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func newTwitterClient(token accessToken) client {
	return client{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}

func (c *client) initTwitterClient() *twitter.Client {

	client := &twitter.Client{
		Authorizer: authorize{
			Token: c.AccessToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	// go c.refreshTwitterClient()
	return client
}
