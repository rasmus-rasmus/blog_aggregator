package handlers

import (
	"blog_aggregator/internal/database"
	"blog_aggregator/utils"
	"strings"

	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (conf *ApiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := struct {
		Name string `json:"name"`
	}{}
	decoderErr := decoder.Decode(&reqBody)
	if decoderErr != nil {
		utils.RespondWithError(w, 500, decoderErr.Error())
		return
	}
	user, createUserErr := conf.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqBody.Name,
	})
	if createUserErr != nil {
		utils.RespondWithError(w, 500, createUserErr.Error())
		return
	}
	utils.RespondWithJSON(w, 201, user)
}

func (conf *ApiConfig) GetUserByKeyHandler(w http.ResponseWriter, r *http.Request) {
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
	utils.RespondWithJSON(w, 200, user)
}
