package main

import (
	"blog_aggregator/internal/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	responseBody := struct {
		Status string `json:"status"`
	}{"ok"}
	respondWithJSON(w, 200, responseBody)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}

func (conf *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := struct {
		Name string `json:"name"`
	}{}
	decoderErr := decoder.Decode(&reqBody)
	if decoderErr != nil {
		respondWithError(w, 500, decoderErr.Error())
		return
	}
	ctx := context.Background()
	user, createUserErr := conf.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqBody.Name,
	})
	if createUserErr != nil {
		respondWithError(w, 500, decoderErr.Error())
		return
	}
	respondWithJSON(w, 201, user)
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	databaseURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	conf := apiConfig{database.New(db)}

	mainRouter := chi.NewRouter()
	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	apiRouter := chi.NewRouter()

	apiRouter.Get("/readiness", readinessHandler)
	apiRouter.Get("/err", errorHandler)
	apiRouter.Post("/users", conf.createUserHandler)

	mainRouter.Mount("/v1/", apiRouter)

	srv := &http.Server{Addr: ":" + port, Handler: mainRouter}

	fmt.Printf("Server listening on port: %s\n", port)

	log.Fatal(srv.ListenAndServe())
}
