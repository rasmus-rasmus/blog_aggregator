package handlers

import (
	"blog_aggregator/internal/database"
	"blog_aggregator/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (conf *ApiConfig) HandlerFeedsPost(w http.ResponseWriter, r *http.Request, user database.User) {
	reqBody := struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}{}
	reqBody, decoderErr := utils.DecodeRequestBody(r, reqBody)
	if decoderErr != nil {
		utils.RespondWithError(w, 500, decoderErr.Error())
		return
	}

	feed, createFeedErr := conf.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqBody.Name,
		Url:       reqBody.Url,
		UserID:    user.ID,
	})
	if createFeedErr != nil {
		utils.RespondWithError(w, 500, createFeedErr.Error())
		return
	}

	follow, createFollowError := conf.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if createFollowError != nil {
		utils.RespondWithError(w, 500, createFollowError.Error())
		return
	}

	responseBody := struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}{databaseFeedToFeed(feed), databaseFollowToFollow(follow)}

	utils.RespondWithJSON(w, 201, responseBody)
}

func (conf *ApiConfig) HandlerFeedsGet(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := conf.DB.GetFeeds(r.Context())
	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}
	feeds := make([]Feed, 0, len(dbFeeds))
	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, databaseFeedToFeed(dbFeed))
	}
	utils.RespondWithJSON(w, 200, feeds)
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
	var dbFeedTime *time.Time
	if dbFeed.LastFetchedAt.Valid {
		dbFeedTime = &dbFeed.LastFetchedAt.Time
	}
	return Feed{dbFeed.ID, dbFeed.CreatedAt, dbFeed.UpdatedAt, dbFeed.Name, dbFeed.Url, dbFeed.UserID, dbFeedTime}
}
