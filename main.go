package main

import (
	"blog_aggregator/handlers"
	"blog_aggregator/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	databaseURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	conf := handlers.ApiConfig{DB: database.New(db)}

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

	apiRouter.Get("/readiness", handlers.ReadinessHandler)
	apiRouter.Get("/err", handlers.ErrorHandler)
	apiRouter.Post("/users", conf.CreateUserHandler)
	apiRouter.Get("/users", conf.GetUserByKeyHandler)

	mainRouter.Mount("/v1/", apiRouter)

	srv := &http.Server{Addr: ":" + port, Handler: mainRouter}

	fmt.Printf("Server listening on port: %s\n", port)

	log.Fatal(srv.ListenAndServe())
}
