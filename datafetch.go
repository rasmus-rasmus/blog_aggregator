package main

import (
	"blog_aggregator/handlers"
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
	Guid    string `xml:"guid"`
	Desc    string `xml:"description"`
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
		log.Printf("Fetching next %d feeds.\n", len(feedsToFetch))
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

				log.Printf("Fetched data from %s at: %s", rssFeed.Title, dbFeed.Url)
			}()
			wg.Wait()
		}
		log.Printf("Succesfully fetched data from %d feeds.", len(feedsToFetch))
		// time.Sleep(intervalBetweenFetches)
	}
}
