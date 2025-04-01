package main

import (
	"database/sql"
	"os"
	"sync/atomic"
	"github.com/evanlin0514/Chirpy/internal/database"
    "github.com/joho/godotenv"
)

type apiConfig struct {
	db *database.Queries
	fileserverHits atomic.Int32
	platform string
}

func NewApiConfig(db *sql.DB) (*apiConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	
	dbQueries := database.New(db)
	return &apiConfig{
		db: dbQueries,
		platform: os.Getenv("PLATFORM"),
	}, nil
}