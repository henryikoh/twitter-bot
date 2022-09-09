package main

import (
	"fmt"
	"time"
)

func main() {

	link, state := newOAuth2link()
	fmt.Printf("url : %s, state %s \n", link, state)

	// call back server listin for request and verifying the state
	runServer(state)

	// jobs.createFilterStream(twitter.TweetSearchStreamRule{
	// 	Value: `("Peter Obi" OR Baba Ahmed Datti) -is:quote -is:retweet -is:reply (hope OR win OR vote OR happy OR excited OR elated OR favorite OR fav OR amazing OR voting) -IPOB -lose -"not win" -"no win" -lie -"wont win" -#BBNaijaS7 -#BBNaija - #Phyna (#PeterObiForPresident OR #Obidatti023 OR #PeterObi4President2023  OR has:hashtags)`,
	// 	Tag:   "Obi4president",
	// })

	// jobs.createFilterStream(twitter.TweetSearchStreamRule{
	// 	Value: "from:PeterObi -is:retweet -is:reply",
	// 	Tag:   "presidentweets",
	// })

	// jobs.deleteStreamRule([]twitter.TweetSearchStreamRuleID{"1566189295282552837"})
	// jobs.searchRule("")

	// jobs.deleteStreamRule([]twitter.TweetSearchStreamRuleID{"1566102889868713988"})

	// jobs.sendTweet("Hello this is a test 2202")

	time.Sleep(time.Second * 5)

}
