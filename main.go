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

func run() error {
	_ = godotenv.Load()

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_TOKEN_SECRET"))
	defer api.Close()

	v := url.Values{}
	stream := api.UserStream(v)

	userIDs := make(map[int64]struct{})
	for _, part := range strings.Split(os.Getenv("TARGET_USER_IDS"), " ") {
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return err
		}
		userIDs[id] = struct{}{}
	}

	for item := range stream.C {
		switch status := item.(type) {
		case anaconda.Tweet:
			if status.Retweeted {
				break
			}
			if _, ok := userIDs[status.User.Id]; !ok {
				break
			}
			if status.InReplyToUserID != 0 && status.InReplyToUserID != status.User.Id {
				break
			}
			if _, err := api.Retweet(status.Id, false); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
