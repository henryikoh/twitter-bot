package main

type TwitterUser struct {
}

type MyJsonName struct {
	One struct {
		PinnedTweet interface{} `json:"PinnedTweet"`
		User        struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"User"`
	} `json:"gg"`
}

type User struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
}
