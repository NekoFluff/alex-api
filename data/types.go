package data

import (
	"time"
)

type ChannelFeed struct {
	FirstName  string
	LastName   string
	TopicURL   string
	Group      string
	Generation int
}

type Livestream struct {
	Author  string    `bson:"author"`
	Url     string    `bson:"url"`
	Date    time.Time `bson:"date"`
	Title   string    `bson:"title"`
	Updated time.Time `bson:"updated"`
}

type Subscription struct {
	User         string `bson:"user"`
	Subscription string `bson:"subscription"`
}

type SubscriptionGroup struct {
	User  string `bson:"_id"`
	Count int    `bson:"count"`
}

type TwitterAuthor struct {
	Id       string `bson:"id"`
	Name     string `bson:"name"`
	UserName string `bson:"username"`
}

type TwitterMedia struct {
	Author            TwitterAuthor `bson:"author"`
	TweetId           string        `bson:"tweet_id"`
	Url               string        `bson:"url"`
	Updated           time.Time     `bson:"updated"`
	CreatedAt         time.Time     `bson:"created_at"`
	PossiblySensitive bool          `bson:"possibly_sensitive"`
	Width             int16         `bson:"width"`
	Height            int16         `bson:"height"`
}
