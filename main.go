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
	log.SetFlags(log.Lshortfile)

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
			log.Fatal("Error parsing TARGET_USER_IDS", err)
		}
		userIds[id] = struct{}{}
	}

	for item := range stream.C {
		switch status := item.(type) {
		case anaconda.Tweet:
			if status.Retweeted {
				break
			}
			if _, ok := userIds[status.User.Id]; !ok {
				break
			}
			if _, err := api.Retweet(status.Id, false); err != nil {
				log.Println(err)
			}
		default:
		}
	}
}
