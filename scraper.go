package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ronkiker/playingwithsql/blob/dev/internal/database"
)

func startScraper(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration) {
	log.Printf("Scrape started on %v goroutines every %s duration \n", concurrency, timeBetweenRequest)

	// how to make the requests on this interval
	ticker := time.NewTicker(timeBetweenRequest)
	// every time a new value comes across the ticker, iterate
	// C is a channel
	for ; ; <-ticker.C {
		log.Println("Scrape started")
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("Error getting feeds to fetch: ", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.SetFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error setting feed as fetched: ", err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error getting feed: ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		//parse date from string
		pubbed, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println("Error parsing time: ", err)
			continue
		}
		_, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				Url:         item.Link,
				PublishedAt: pubbed,
				FeedID:      feed.ID,
			})
		if err != nil {
			log.Println("Error creating post: ", err)
			return
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
