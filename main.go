package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main(){
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error connecting to database.")
	}
	config, err := NewApiConfig(db)
	if err != nil{
		log.Fatal("error creating config:", err)
	}

	mux := http.NewServeMux()
	homepage := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", config.middlewareMetricsInc(homepage))
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("GET /admin/metrics", config.showCounter)
	mux.HandleFunc("POST /admin/reset", config.resetUser)
	mux.HandleFunc("GET /api/healthz", handler)
	mux.HandleFunc("POST /api/validate_chirp", validHanlder)
	mux.HandleFunc("POST /api/users", config.addUser)
	
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	
	if err := server.ListenAndServe(); err != nil{
		fmt.Println(err)
	}
}