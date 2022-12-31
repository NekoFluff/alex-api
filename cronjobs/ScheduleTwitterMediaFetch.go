package cronjobs

import (
	"alex-api/data"
	"alex-api/internal/twitterapi"
	"context"

	"log"
	"time"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

func ScheduleTwitterMediaFetch() *cron.Cron {
	go getInArtTwitterMedia()

	c := cron.New()
	spec := "0 */5 * * * *"
	c.AddFunc(spec, func() {
		getInArtTwitterMedia()
	})
	c.Start()

	log.Printf("Scheduled Twitter Media fetch for every 5 minutes (%v)\n", spec)
	return c
}

func getInArtTwitterMedia() {
	query := `(#inART OR #いなート) has:media -is:retweet`

	opts := twitter.TweetRecentSearchOpts{
		MediaFields: []twitter.MediaField{twitter.MediaFieldURL, twitter.MediaFieldType, twitter.MediaFieldPublicMetrics, twitter.MediaFieldMediaKey, twitter.MediaFieldHeight, twitter.MediaFieldWidth},
		Expansions:  []twitter.Expansion{twitter.ExpansionEntitiesMentionsUserName, twitter.ExpansionAuthorID, twitter.ExpansionAttachmentsMediaKeys},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldConversationID, twitter.TweetFieldAttachments, twitter.TweetFieldSource, twitter.TweetFieldAuthorID, twitter.TweetFieldPossiblySensitve},
	}

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	db := data.New(log.WithContext(context.TODO()))

	twitterMedia, err := db.GetMostRecentTwitterMedia()
	if err != nil {
		log.Println(err)
		return
	}

	weekAgo := time.Now().Add(-7 * 24 * time.Hour)
	if twitterMedia.CreatedAt.After(weekAgo) {
		opts.SinceID = twitterMedia.TweetId
	} else {
		opts.StartTime = weekAgo
	}
	twitterapi.GetTwitterMedia(query, opts)

}
