package handlers

import (
	"blog_aggregator/utils"
	"net/http"
)

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	responseBody := struct {
		Status string `json:"status"`
	}{"ok"}
	utils.RespondWithJSON(w, 200, responseBody)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, 500, "Internal Server Error")
}
