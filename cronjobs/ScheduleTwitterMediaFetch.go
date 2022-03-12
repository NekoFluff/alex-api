package cronjobs

import (
	"addi/data"
	"addi/twitterapi"

	"log"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/robfig/cron"
)

func ScheduleTwitterMediaFetch() *cron.Cron {

	c := cron.New()
	spec := "0 */5 * * * *"
	c.AddFunc(spec, func() {
		query := `"#inART" has:media -is:retweet`

		opts := twitter.TweetRecentSearchOpts{
			MediaFields: []twitter.MediaField{twitter.MediaFieldURL, twitter.MediaFieldType, twitter.MediaFieldPublicMetrics, twitter.MediaFieldMediaKey, twitter.MediaFieldHeight, twitter.MediaFieldWidth},
			Expansions:  []twitter.Expansion{twitter.ExpansionEntitiesMentionsUserName, twitter.ExpansionAuthorID, twitter.ExpansionAttachmentsMediaKeys},
			TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldConversationID, twitter.TweetFieldAttachments, twitter.TweetFieldSource, twitter.TweetFieldAuthorID, twitter.TweetFieldPossiblySensitve},
		}

		twitterMedia, err := data.GetMostRecentTwitterMedia()
		if err != nil {
			log.Println(err)
		} else {
			opts.SinceID = twitterMedia.TweetId
			twitterapi.GetTwitterMedia(query, opts)
		}
	})
	c.Start()

	log.Printf("Scheduled Twitter Media fetch for every 5 minutes (%v)\n", spec)
	return c
}
