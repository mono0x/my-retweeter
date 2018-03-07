package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

func run() error {
	_ = godotenv.Load()

	api := anaconda.NewTwitterApiWithCredentials(
		os.Getenv("TWITTER_OAUTH_TOKEN"),
		os.Getenv("TWITTER_OAUTH_TOKEN_SECRET"),
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"))
	defer api.Close()

	userIDStrs := strings.Split(os.Getenv("TARGET_USER_IDS"), " ")
	userIDs := make([]int64, 0, len(userIDStrs))
	sinceIDs := make(map[int64]int64, len(userIDStrs))
	for _, part := range userIDStrs {
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return err
		}
		userIDs = append(userIDs, id)
		sinceIDs[id] = 0
	}

	interval := time.Duration(float64(15*time.Minute) / (180.0 / float64(len(userIDs))))
	if interval < 1*time.Minute {
		interval = 1 * time.Minute
	}
	t := time.NewTicker(interval)
	defer t.Stop()

	for _ = range t.C {
		for _, userID := range userIDs {
			var v url.Values
			v.Set("count", "200")
			v.Set("exclude_replies", "true")
			sinceID := sinceIDs[userID]
			if sinceID > 0 {
				v.Set("since_id", strconv.FormatInt(sinceID, 10))
			}
			timeline, err := api.GetUserTimeline(v)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, status := range timeline {
				if status.Id > sinceID {
					sinceID = status.Id
				}
				if status.Retweeted {
					continue
				}
				if _, err := api.Retweet(status.Id, false); err != nil {
					log.Println(err)
				}
			}
			sinceIDs[userID] = sinceID
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
