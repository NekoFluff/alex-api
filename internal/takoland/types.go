package takoland

type TwitterMedia struct {
	Author            *string `json:"author"`
	CreatedAt         *string `json:"created_at"`
	Height            *int64  `json:"height"`
	PossiblySensitive *bool   `json:"possiblySensitive"`
	TweetID           *string `json:"tweetId"`
	URL               *string `json:"url"`
	Width             *int64  `json:"width"`
}
