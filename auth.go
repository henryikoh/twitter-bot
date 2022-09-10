package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	auth2Url      = "https://twitter.com/i/oauth2/authorize"
	auth2TokenUrl = "https://api.twitter.com/2/oauth2/token"
	client_id     = "c0F3WVVSa04zWl9CNDIyLURGT2c6MTpjaQ"
	client_secret = "xnti0tH34WW2Ez24U4rMXwRHO1j-hCAJFW-QkzKSiK6SoiIB_4"
	redirect_uri  = "https://twitter-bot-w7ofesn2oa-uc.a.run.app:8080/twitback"
	state         = "xyzABC1235"
	state2        = "xyzABC1234567"
)

type scope struct {
	scope []string
}

// this function allows you added twitter scopes to the Auth handler you can pass in a string eg. "offline./// access" or an array of scops
func (s *scope) AddScopes(elems ...string) {
	s.scope = append(s.scope, elems...)
}

func (s *scope) toString() string {
	if len(s.scope) == 0 {
		fmt.Println("scope cant be empty")
	}
	return strings.Join(s.scope, " ")
}

func (s *scope) espacePath() string {
	return url.PathEscape(s.toString())
}

type Auth2ulrBuilder struct {
	url                   string
	scope                 *scope
	client_id             string
	state                 string
	redirect_uri          string
	code_challenge        string
	code_challenge_method string
	response_type         string
}

type accessToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (auth2 *Auth2ulrBuilder) generate() string {

	u, err := url.Parse(auth2.url)
	if err != nil {
		log.Fatal(err)
	}

	parsescope := auth2.scope.espacePath()

	values := u.Query()

	values.Set("response_type", "code")
	values.Add("scope", parsescope)
	values.Add("code_challenge", "challenge")
	values.Add("state", state)
	values.Add("code_challenge_method", "plain")
	values.Set("redirect_uri", redirect_uri)
	values.Add("client_id", client_id)

	decoded, err := url.QueryUnescape(values.Encode())
	if err != nil {
		log.Println(err)
	}
	u.RawQuery = decoded

	return u.String()

}

func getAccesToken(code string) {

	client := &http.Client{}

	value := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirect_uri},
		"code_verifier": {"challenge"},
		"client_id":     {client_id},
	}

	// fmt.Println(value.Encode())
	// fmt.Println(strings.NewReader(value.Encode()))
	strin, err := url.QueryUnescape(value.Encode())
	if err != nil {
		log.Println(err)
	}

	// fmt.Printf("String: %s \n", strin)

	req, _ := http.NewRequest("POST", auth2TokenUrl, strings.NewReader(strin))

	req.SetBasicAuth(client_id, client_secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	result := accessToken{}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}
	// sb := string(body)

	// store access and refresh tokens.
	log.Println(result)
	db := InitDAO()

	res := db.db.Create(&result)
	if res.Error != nil {
		log.Fatal(result)
	}

	// run jobs
	jobs := initJobs()

	jobs.consumeStream()

}
func refreshAccesToken(refreshToken string) accessToken {
	fmt.Println("refresh token hit")
	// code to get a new access token from the refresh token
	client := &http.Client{}

	value := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	req, _ := http.NewRequest("POST", auth2TokenUrl, strings.NewReader(value.Encode()))

	req.SetBasicAuth(client_id, client_secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))

	result := accessToken{}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	fmt.Printf("new refresh token %s", result)

	db := InitDAO()

	res := db.db.Model(&result).Where("refresh_token = ?", refreshToken).Save(result)

	fmt.Printf("rows affected %s", res.RowsAffected)

	if res.Error != nil {
		log.Fatal(result)
	}

	return result
}
func newOAuth2link() (link string, stat string) {

	// // create scope variable slice
	// gg := &scope{}

	// // gg.AddScopes("offline.access", "tweet.read")

	url := Auth2ulrBuilder{
		url: auth2Url,
		scope: &scope{
			scope: []string{"offline.access", "tweet.read", "tweet.write", "users.read", "follows.write", "like.write"},
		},
		client_id:             client_id,
		state:                 state,
		redirect_uri:          redirect_uri,
		code_challenge:        state,
		code_challenge_method: "plain",
		response_type:         "code",
	}

	link = url.generate()

	stat = url.state

	return

	// Use the Query() method to get the query string params as a url.Values map.

}
