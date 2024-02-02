package main

import (
	"blog_aggregator/handlers"
	"blog_aggregator/internal/database"
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
	Guid    string `xml:"guid"`
	Desc    string `xml:"description"`
}

func parseTimeString(timestamp string) (time.Time, bool) {
	parsedTime, err := time.Parse(time.RFC1123Z, timestamp)
	if err == nil {
		return parsedTime, true
	}
	parsedTime, err = time.Parse(time.RFC1123, timestamp)
	if err == nil {
		return parsedTime, true
	}
	return time.Time{}, false
}

func writePostToDB(item Item, feedId uuid.UUID, conf *handlers.ApiConfig) (database.Post, error) {
	parsedTime, success := parseTimeString(item.PubDate)
	if !success {
		return database.Post{}, errors.New(fmt.Sprintf("Couldn't parse pubDate: %s", item.PubDate))
	}
	createPostParams := database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Title:       sql.NullString{String: item.Title, Valid: len(item.Title) > 0},
		Url:         item.Link,
		Description: sql.NullString{String: item.Desc, Valid: len(item.Desc) > 0},
		PublishedAt: parsedTime,
		FeedID:      feedId,
	}
	return conf.DB.CreatePost(context.Background(), createPostParams)
}

type RssFeed struct {
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`
	Items       []Item `xml:"channel>item"`
}

func fetchDataFromFeed(url string) (RssFeed, error) {
	feed := RssFeed{}

	res, getErr := http.Get(url)
	if getErr != nil {
		return feed, getErr
	}
	resBody, readErr := io.ReadAll(res.Body)
	res.Body.Close()
	if readErr != nil {
		return feed, readErr
	}
	unmarshalErr := xml.Unmarshal(resBody, &feed)
	if unmarshalErr != nil {
		return feed, unmarshalErr
	}
	return feed, nil
}

func dataFetchingWorker(numFeedsToFetch int, intervalBetweenFetches time.Duration, conf *handlers.ApiConfig) {
	ticker := time.NewTicker(intervalBetweenFetches)
	for ; ; <-ticker.C {
		feedsToFetch, nextFeedsErr := conf.DB.GetNextFeedsToFetch(context.Background(), int32(numFeedsToFetch))
		if nextFeedsErr != nil {
			log.Println("Couldn't get next feeds to fetch.")
		}
		formatLoggerHeader(fmt.Sprintf("Fetching next %d feeds.", len(feedsToFetch)))
		var wg sync.WaitGroup
		for _, dbFeed := range feedsToFetch {
			wg.Add(1)
			go func() {
				defer wg.Done()
				rssFeed, fetchErr := fetchDataFromFeed(dbFeed.Url)
				if fetchErr != nil {
					log.Printf("Couldn't fetch data from %s\n", dbFeed.Url)
					return
				}
				if _, markErr := conf.DB.MarkFeedFetched(context.Background(), dbFeed.ID); markErr != nil {
					log.Printf("Couldn't update fetched feed: %s\n", dbFeed.Url)
				}

				newPostsCounter := 0
				for _, item := range rssFeed.Items {
					post, writeErr := writePostToDB(item, dbFeed.ID, conf)
					if writeErr != nil {
						if writeErr.Error() == `pq: duplicate key value violates unique constraint "posts_url_key"` {
							continue
						}
						log.Printf("Couldn't write post: %s. %s", item.Link, writeErr.Error())
					} else {
						log.Printf("Wrote post to DB: %s\n", post.Url)
						newPostsCounter++
					}
				}
				log.Printf("Scraped %d new posts from %s\n", newPostsCounter, dbFeed.Url)
			}()
			wg.Wait()
		}
		formatLoggerFooter(fmt.Sprintf("Succesfully fetched data from %d feeds.", len(feedsToFetch)))
	}
}

func formatLoggerHeader(textToLog string) {
	fmt.Println(strings.Repeat("~", len(textToLog)+20))
	log.Println(textToLog)
}

func formatLoggerFooter(textToLog string) {
	log.Println(textToLog)
	fmt.Println(strings.Repeat("~", len(textToLog)+20))
}
