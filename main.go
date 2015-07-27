package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
	"reflect"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_TOKEN_SECRET"))
	api.SetLogger(anaconda.BasicLogger)

	v := url.Values{}
	stream := api.UserStream(v)

	targetUserId, err := strconv.ParseInt(os.Getenv("TARGET_USER_ID"), 10, 64)
	if err != nil {
		log.Fatal("Error parsing TARGET_USER_ID")
	}

loop:
	for {
		select {
		case item := <-stream.C:
			fmt.Println(reflect.TypeOf(item))
			switch status := item.(type) {
			case anaconda.Tweet:
				if !status.Retweeted && status.User.Id == targetUserId {
					_, _ = api.Retweet(status.Id, false)
				}
			default:
			}
		case <-stream.Quit:
			break loop
		}
	}
}
