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
