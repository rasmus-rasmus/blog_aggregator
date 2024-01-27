package handlers

import (
	"blog_aggregator/internal/database"
	"blog_aggregator/utils"

	"net/http"
	"time"

	"github.com/google/uuid"
)

func (conf *ApiConfig) HandlerUsersPost(w http.ResponseWriter, r *http.Request) {
	reqBody := struct {
		Name string `json:"name"`
	}{}
	reqBody, decoderErr := utils.DecodeRequestBody(r, reqBody)
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

func HandlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	utils.RespondWithJSON(w, 200, user)
}
