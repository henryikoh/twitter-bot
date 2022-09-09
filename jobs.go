package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	twitter "github.com/g8rswimmer/go-twitter/v2"
	"gorm.io/gorm"
)

type jobs struct {
	client    *twitter.Client
	id        string
	db        *gorm.DB
	token     accessToken
	appToken  string
	rtalike   chan *twitter.TweetMessage
	vidlink   chan *twitter.TweetMessage
	president chan *twitter.TweetMessage
}

func initJobs() *jobs {

	var token accessToken
	db := InitDAO()

	// the current token from the database access token might have already expeired.
	db.db.Find(&token)
	fmt.Printf("token from db %s", token)

	// call refresh token on init incase the token has expired
	newtoken := refreshAccesToken(token.RefreshToken)

	// new token is passed into the twitter client
	c := newTwitterClient(newtoken)
	twitterClient := c.initTwitterClient()

	// twitterClient.AuthUserLookup()
	user := userLookupme(*twitterClient)

	jobs := &jobs{
		client:    c.initTwitterClient(),
		id:        user.ID,
		db:        db.db,
		token:     newtoken,
		appToken:  "AAAAAAAAAAAAAAAAAAAAANQHgQEAAAAA6AOvjbW2hIoxyeit0ectjVfkTgo%3DeTTtbPqq9Zu8R0NHVAwdz561yymC88dTeKOHjFRqHxv456nQBm",
		rtalike:   make(chan *twitter.TweetMessage, 50),
		vidlink:   make(chan *twitter.TweetMessage, 50),
		president: make(chan *twitter.TweetMessage, 50),
	}
	go jobs.refreshTwitterClient()
	return jobs
}

// this runs ever 1hour 30mins and it gets a new refresh token for the twitter client
func (j *jobs) refreshTwitterClient() {

	ticker := time.NewTicker(time.Second * 5400)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		// this sends a ping to the connect very 54 seconds
		case <-ticker.C:
			fmt.Println("ticker hit")
			token := refreshAccesToken(j.token.RefreshToken)
			j.token = token
			c := newTwitterClient(token)
			j.client = c.initTwitterClient()
		case tweets := <-j.rtalike:
			for _, tweet := range tweets.Raw.Tweets {
				j.retweetandlike(tweet.ID)
			}
		case <-j.vidlink:
			// qoute tweet with download link
			j.consumeStream()
		case tweets := <-j.president:
			for _, tweet := range tweets.Raw.Tweets {
				j.preisdentTweet(tweet)
			}
			// this channel is used to like and retweet tweets
		}
	}

}

/**
	In order to run, the user will need to provide the bearer token and the list of tweet ids.
**/
func (j *jobs) sendTweet(text string) {
	j.client.CreateTweet(context.Background(), twitter.CreateTweetRequest{
		Text: text,
	})

}
func (j *jobs) createFilterStream(streamRule twitter.TweetSearchStreamRule) {

	client := &twitter.Client{
		Authorizer: authorize{
			Token: j.appToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	fmt.Println("Callout to tweet search stream add rule callout")

	searchStreamRules, err := client.TweetSearchStreamAddRule(context.Background(), []twitter.TweetSearchStreamRule{streamRule}, false)
	if err != nil {
		log.Panicf("tweet search stream add rule callout error: %v", err)
	}

	enc, err := json.MarshalIndent(searchStreamRules, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))
}

func (j *jobs) searchRule() {

	client := &twitter.Client{
		Authorizer: authorize{
			Token: j.appToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	fmt.Println("Callout to tweet search stream rules callout")

	searchStreamRules, err := client.TweetSearchStreamRules(context.Background(), []twitter.TweetSearchStreamRuleID{})
	if err != nil {
		log.Panicf("tweet search stream rule callout error: %v", err)
	}

	enc, err := json.MarshalIndent(searchStreamRules, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))

}
func (j *jobs) deleteStreamRule(rulesID []twitter.TweetSearchStreamRuleID) {
	client := &twitter.Client{
		Authorizer: authorize{
			Token: j.appToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	fmt.Println("Callout to tweet search stream delete rule callout")

	ruleIDs := rulesID

	searchStreamRules, err := client.TweetSearchStreamDeleteRuleByID(context.Background(), ruleIDs, false)
	if err != nil {
		log.Panicf("tweet search stream delete rule callout error: %v", err)
	}

	enc, err := json.MarshalIndent(searchStreamRules, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))
}

func (j *jobs) consumeStream() {

	client := &twitter.Client{
		Authorizer: authorize{
			Token: j.appToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	opts := twitter.TweetSearchStreamOpts{
		Expansions:  []twitter.Expansion{twitter.ExpansionAttachmentsMediaKeys, twitter.ExpansionAuthorID},
		MediaFields: []twitter.MediaField{twitter.MediaFieldURL},
		TweetFields: []twitter.TweetField{twitter.TweetFieldAttachments},
		UserFields:  []twitter.UserField{twitter.UserFieldUserName, twitter.UserFieldProfileImageURL},
	}

	fmt.Println("Callout to tweet search stream callout")

	outputFile, err := os.Create("test")
	if err != nil {
		log.Panicf("tweet stream output file error %v", err)
	}
	defer outputFile.Close()

	tweetStream, err := client.TweetSearchStream(context.Background(), opts)

	if err != nil {
		log.Panicf("tweet sample callout error: %v", err)
	}
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	func() {
		defer tweetStream.Close()
		for {
			select {
			case <-ch:
				fmt.Println("closing")
				return
			case tm := <-tweetStream.Tweets():
				// send tweet to write and retweet channel..
				j.rtalike <- tm

				user := tm.Raw.Includes.Users

				// presidential tweet
				if user[0].UserName == "PeterObi" {
					// presidential tweet
					j.president <- tm
				}

				fmt.Println(tm.Raw.Tweets[0])
				fmt.Println(tm.Raw)
				tmb, err := json.Marshal(tm)
				if err != nil {
					fmt.Printf("error decoding tweet message %v", err)
				}
				outputFile.WriteString(fmt.Sprintf("tweet: %s\n\n", string(tmb)))
				outputFile.Sync()
				fmt.Printf("new tweet,%v", tm)
			case sm := <-tweetStream.SystemMessages():
				smb, err := json.Marshal(sm)
				if err != nil {
					fmt.Printf("error decoding system message %v", err)
				}
				outputFile.WriteString(fmt.Sprintf("system: %s\n\n", string(smb)))
				outputFile.Sync()
				fmt.Println("system")
			case strErr := <-tweetStream.Err():
				outputFile.WriteString(fmt.Sprintf("error: %v\n\n", strErr))
				outputFile.Sync()
				fmt.Println("error")
			}
			if tweetStream.Connection() == false {
				fmt.Println("connection lost")
				// reconnect to stream
				// j.vidlink <-
				return
			}
		}
	}()
}

func userLookupme(client twitter.Client) twitter.UserObj {
	opts := twitter.UserLookupOpts{
		Expansions: []twitter.Expansion{twitter.ExpansionPinnedTweetID},
	}
	fmt.Println("Callout to auth user lookup callout")

	userResponse, err := client.AuthUserLookup(context.Background(), opts)
	if err != nil {
		log.Panicf("auth user lookup error: %v", err)
	}

	result := userResponse.Raw.Users
	fmt.Println(len(result))

	fmt.Println("one result found")
	return *result[0]

	// fmt.Println(string(enc))
}

func (j *jobs) retweetandlike(id string) {

	_, err := j.client.UserLikes(context.Background(), j.id, id)
	if err != nil {
		log.Panicf("user like tweet error: %v", err)
	}
	fmt.Printf("tweet %s liked", id)

	//
	_, err2 := j.client.UserRetweet(context.Background(), j.id, id)
	if err != nil {
		log.Panicf("user retweet error: %v", err2)
	}

	fmt.Printf("tweet %s retweeted", id)

}

func (j *jobs) preisdentTweet(tweet *twitter.TweetObj) {

	rand.Seed(time.Now().UnixNano())

	qoutes := []string{"Our president has given us a word \n#PeterObiForPresident", "My presient \n#PeterObiForPresident", "Ride on Sir", "Are you Obidient and Yusful?", "The man of the people", "Incoming>>> \n#PeterObiForPresident", "OBI OBI OBI \n\n#PeterObiForPresident"}

	replies := []string{"Fellow obidients kindly follow me, I am a bot built by @henryikoh_", "Even robots are obidient now... follow for lastest updates currated from the community", "Nah OBI I wan dey follow now", "Nigeria must be great!!!!", "OBI\nOBI\nOBI\nOBI\nOBI\nOBI\n#PeterObiForPresident", "Obidient gather here let’s follow each other", "When I grow up I would love to work for mr OBI... sorry @henryikoh_", "Ride on Sir", "My code gets updated everyday to make me more more Obidient and Yusful... follow for curated news from the obidient family ❤️"}

	j.replyTweet(tweet.ID, replies[rand.Intn(len(replies))])
	j.qouteTweet(tweet.ID, qoutes[rand.Intn(len(qoutes))])

	j.retweetandlike(tweet.ID)

}

func (j *jobs) qouteTweet(id string, text string) {
	req := twitter.CreateTweetRequest{
		Text:         text,
		QuoteTweetID: id,
	}
	fmt.Println("Callout to create tweet callout")

	_, err := j.client.CreateTweet(context.Background(), req)
	if err != nil {
		log.Panicf("create tweet error: %v", err)
	}

}

func (j *jobs) replyTweet(id string, text string) {
	req := twitter.CreateTweetRequest{
		Text: text,
		Reply: &twitter.CreateTweetReply{
			InReplyToTweetID: id,
		},
	}
	fmt.Println("Callout to create tweet callout")

	_, err := j.client.CreateTweet(context.Background(), req)
	if err != nil {
		log.Panicf("create tweet error: %v", err)
	}
}
