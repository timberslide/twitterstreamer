package main

import (
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/kelseyhightower/envconfig"
	"github.com/timberslide/gotimberslide"
)

type specification struct {
	TsToken        string `envconfig:"ts_token"`
	TsTopic        string `envconfig:"ts_topic"`
	ConsumerKey    string `envconfig:"consumer_key"`
	ConsumerSecret string `envconfig:"consumer_secret"`
	AccessToken    string `envconfig:"access_token"`
	AccessSecret   string `envconfig:"access_secret"`
}

func main() {
	var err error

	var s specification
	err = envconfig.Process("APP", &s)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s)

	// Configure our http client
	config := oauth1.NewConfig(s.ConsumerKey, s.ConsumerSecret)
	token := oauth1.NewToken(s.AccessToken, s.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Create a Twitter client
	client := twitter.NewClient(httpClient)

	// Hardcode some words for now
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"Clinton", "Trump", "election", "whitehouse"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to timberslide
	// Create a Timberslide client
	tsClient, err := ts.NewClient("gw.timberslide.com:443", s.TsToken)

	// Connect to Timberslide
	err = tsClient.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer tsClient.Close()

	// Send some messages into a Timberslide topic
	ch, err := tsClient.CreateChannel(s.TsTopic)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Starting stream")
	for message := range stream.Messages {
		switch message := message.(type) {
		case *twitter.Tweet:
			tweet := fmt.Sprintf("%s - @%s - %s", message.CreatedAt, message.User.ScreenName, message.Text)
			ch.Send(tweet)
		case *twitter.StallWarning:
			log.Println("Stall warning!")
		}
	}
	log.Println("Stream exited")
}
