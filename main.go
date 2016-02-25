package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_TOKEN_SECRET"))

	v := url.Values{}
	stream := api.UserStream(v)

	userIds := make(map[int64]struct{})
	for _, part := range strings.Split(os.Getenv("TARGET_USER_IDS"), " ") {
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			log.Fatal("Error parsing TARGET_USER_IDS")
		}
		userIds[id] = struct{}{}
	}

	for {
		select {
		case item := <-stream.C:
			switch status := item.(type) {
			case anaconda.Tweet:
				if !status.Retweeted {
					if _, ok := userIds[status.User.Id]; ok {
						_, _ = api.Retweet(status.Id, false)
					}
				}
			default:
			}
		}
	}
}
