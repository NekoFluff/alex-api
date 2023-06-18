package twitterapi

import (
	"alex-api/internal/data"
	"alex-api/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/sirupsen/logrus"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func GetTwitterMedia(query string, opts twitter.TweetRecentSearchOpts) {
	client := &twitter.Client{
		Authorizer: authorize{
			Token: utils.GetEnvVar("TWITTER_BEARER_TOKEN"),
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.Println("Callout to tweet recent search callout")
	log.Printf("SinceId %v\n", opts.SinceID)

	tweetResponse, err := client.TweetRecentSearch(context.Background(), query, opts)
	if err != nil {
		log.WithFields(logrus.Fields{
			"query": query,
		}).Errorf("tweet lookup error: %v", err)
	}

	dictionaries := tweetResponse.Raw.TweetDictionaries()

	// enc, err := json.MarshalIndent(dictionaries, "", "    ")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// fmt.Println(string(enc))

	metaBytes, err := json.MarshalIndent(tweetResponse.Meta, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(metaBytes))

	for _, tweetDictionary := range dictionaries {
		for _, media := range tweetDictionary.AttachmentMedia {
			if media.URL == "" {
				continue
			}

			createdAt, _ := time.Parse(time.RFC3339, tweetDictionary.Tweet.CreatedAt)

			twitterMedia := &data.TwitterMedia{
				Author: data.TwitterAuthor{
					Id:       tweetDictionary.Author.ID,
					Name:     tweetDictionary.Author.Name,
					Username: tweetDictionary.Author.UserName,
				},
				TweetId:           tweetDictionary.Tweet.ID,
				Url:               media.URL,
				CreatedAt:         createdAt,
				PossiblySensitive: tweetDictionary.Tweet.PossiblySensitive,
				Updated:           time.Now(),
				Width:             int16(media.Width),
				Height:            int16(media.Height),
			}

			log := logrus.New()
			log.SetFormatter(&logrus.JSONFormatter{})
			db := data.New(log.WithContext(context.TODO()))
			defer db.Disconnect()
			db.CreateTwitterMedia(*twitterMedia)
		}
	}

	if tweetResponse.Meta.NextToken != "" {
		fmt.Println(tweetResponse.Meta.NextToken)
		GetTwitterMedia(query, twitter.TweetRecentSearchOpts{
			NextToken:   tweetResponse.Meta.NextToken,
			TweetFields: opts.TweetFields,
			Expansions:  opts.Expansions,
			MediaFields: opts.MediaFields,
			MaxResults:  opts.MaxResults,
		})
	}
}
