package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"github.com/evanlin0514/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db *database.Queries
	fileserverHits atomic.Int32
}

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
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

func respondWithJSON (w http.ResponseWriter, code int, payload interface{}) error {
	res, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "apllication/json")
	w.WriteHeader(code)
	w.Write(res)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
    return respondWithJSON(w, code, map[string]string{"error": msg})
}

type cleanBody struct {
	Body string `json:"cleaned_body"`
}
type parameters struct {
	Body string `json:"body"`
}

func cleanInput(str string) cleanBody{
	clean := cleanBody{
		Body: str,
	}

	banMaps := make(map[string]bool)
	banWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, ban := range banWords{
		banMaps[ban] = true
	}

	words := strings.Split(clean.Body, " ")
	for i, word := range words {
		if banMaps[strings.ToLower(word)] {
			words[i] = "****"
		}
	}
	clean.Body = strings.Join(words, " ")
	return clean
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

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	config := &apiConfig{
		db: dbQueries,
	}

	mux := http.NewServeMux()
	homepage := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", config.middlewareMetricsInc(homepage))

	// Option 2: Use http.StripPrefix to remove the duplicated path segment
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request){
		config.showCounter(w, r)
	})

	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request){
		config.resetCounter(w, r)
	})

	mux.HandleFunc("GET /api/healthz", handler)

	mux.HandleFunc("POST /api/validate_chirp", validHanlder)
	
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	
	if err := server.ListenAndServe(); err != nil{
		fmt.Println(err)
	}
}