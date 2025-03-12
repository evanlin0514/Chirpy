package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	times := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(times))
}

func (cfg *apiConfig)resetCounter(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits.Store(0)
}

func main(){
	mux := http.NewServeMux()
	var config apiConfig
	homepage := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", config.middlewareMetricsInc(homepage))
		// Option 2: Use http.StripPrefix to remove the duplicated path segment
		mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request){
		config.showCounter(w, r)
	})
	mux.HandleFunc("POST /reset", func(w http.ResponseWriter, r *http.Request){
		config.resetCounter(w, r)
	})

	mux.HandleFunc("GET /healthz", handler)
	
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	
	if err := server.ListenAndServe(); err != nil{
		fmt.Println(err)
	}
}