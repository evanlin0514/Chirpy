package main

import (
	"net/http"
	"encoding/json"
	"log"
	"time"
	"fmt"
	"context"
    "github.com/evanlin0514/Chirpy/internal/database"
    "github.com/google/uuid"
)

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validHanlder (w http.ResponseWriter, r *http.Request){
	// decode the request body
	decoder := json.NewDecoder(r.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding params: %v", err)
		respondWithError(w, 500, "something went wrong")
		return
	}

	//check if valid
	if len(param.Body) > 200 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	clean := cleanInput(param.Body)
	respondWithJSON(w, 200, clean)
}

func (cfg *apiConfig) addUser (w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := parameters{}
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("error decoding params: %v", err)
		respondWithError(w, 500, "something went wrong")
		return
	}

	params := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email: body.Body,
	}

	user, err := cfg.db.CreateUser(context.Background(), params)
	if err != nil{
		log.Printf("error creating user: %v", err)
		respondWithError(w, 500, "something went wrong")
		return
	}

	respondWithJSON(w, 201, user)
}

func(cfg *apiConfig) resetUser (w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}
	err := cfg.db.ResetUser(context.Background())
	if err != nil {
		log.Printf("error reseting users table")
		respondWithError(w, 500, "something went wrong")
		return
	}
}

func (cfg *apiConfig)showCounter(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	times := fmt.Sprintf(`<html>
		<body>
		  <h1>Welcome, Chirpy Admin</h1>
		  <p>Chirpy has been visited %d times!</p>
		</body>
	  </html>`, cfg.fileserverHits.Load())
	w.Write([]byte(times))
}


func (cfg *apiConfig)resetCounter(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits.Store(0)
}