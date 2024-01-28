package handlers

import (
	"blog_aggregator/internal/database"
	"blog_aggregator/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (conf *ApiConfig) HandlerFeedFollowsPost(w http.ResponseWriter, r *http.Request, user database.User) {
	reqBody := struct {
		FeedId string `json:"feed_id"`
	}{}
	reqBody, decoderErr := utils.DecodeRequestBody(r, reqBody)
	if decoderErr != nil {
		utils.RespondWithError(w, 500, decoderErr.Error())
		return
	}
	feed_follow, createErr := conf.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    uuid.MustParse(reqBody.FeedId),
		UserID:    user.ID,
	})
	if createErr != nil {
		utils.RespondWithError(w, 400, createErr.Error())
		return
	}
	utils.RespondWithJSON(w, 201, feed_follow)
}

func (conf *ApiConfig) HandlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID := chi.URLParam(r, "feedFollowID")
	if len(feedFollowID) != 36 {
		utils.RespondWithError(w, 400, "Invalid id")
		return
	}
	feedFollow, deleteErr := conf.DB.DeleteFeedFollow(r.Context(), uuid.MustParse(feedFollowID))
	if deleteErr != nil {
		utils.RespondWithError(w, 400, deleteErr.Error())
		return
	}
	utils.RespondWithJSON(w, 200, feedFollow)
}

func (conf *ApiConfig) HandlerFeedFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	follows, err := conf.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		utils.RespondWithError(w, 500, err.Error())
		return
	}
	utils.RespondWithJSON(w, 200, follows)
}
