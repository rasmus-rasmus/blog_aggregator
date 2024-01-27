package handlers

import (
	"blog_aggregator/internal/database"
	"blog_aggregator/utils"
	"net/http"
	"strings"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (conf *ApiConfig) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) != 2 || authHeader[0] != "ApiKey" {
			utils.RespondWithError(w, 401, "Missing authorization")
			return
		}
		user, getUserErr := conf.DB.GetUserByApiKey(r.Context(), authHeader[1])
		if getUserErr != nil {
			utils.RespondWithError(w, 404, "Resource not found")
			return
		}
		handler(w, r, user)
	}
}
