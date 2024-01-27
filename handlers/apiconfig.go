package handlers

import "blog_aggregator/internal/database"

type ApiConfig struct {
	DB *database.Queries
}
