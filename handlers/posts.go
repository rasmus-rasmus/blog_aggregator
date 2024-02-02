package handlers

import (
	"blog_aggregator/internal/database"
	"blog_aggregator/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (conf *ApiConfig) HandlerPostsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	limit, atoiErr := strconv.Atoi(r.URL.Query().Get("limit"))
	if atoiErr != nil {
		limit = 10
	}

	dbPosts, getErr := conf.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if getErr != nil {
		utils.RespondWithError(w, 500, getErr.Error())
		return
	}

	posts := make([]Post, 0, len(dbPosts))
	for _, dbPost := range dbPosts {
		posts = append(posts, databasePostToPost(dbPost))
	}

	utils.RespondWithJSON(w, 200, posts)
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       *string   `json:"title"`
	Url         string    `json:"url"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePostToPost(dbPost database.Post) Post {
	var postTitle, postDescription *string
	if dbPost.Title.Valid {
		postTitle = &dbPost.Title.String
	}
	if dbPost.Description.Valid {
		postDescription = &dbPost.Description.String
	}
	return Post{dbPost.ID, dbPost.CreatedAt, dbPost.UpdatedAt, postTitle, dbPost.Url, postDescription, dbPost.PublishedAt, dbPost.FeedID}
}
