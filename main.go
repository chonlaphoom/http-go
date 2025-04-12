package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type HitCounter interface {
	middlewareMetricInt(next http.Handler) http.Handler
	getHits() int32
	resetHits()
}

func (cfg *apiConfig) middlewareMetricInt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getHits() int32 {
	return cfg.fileserverHits.Load()
}

func (cfg *apiConfig) resetHits() {
	cfg.fileserverHits.Store(0)
}

func main() {
	port := "8080"
	address := ":" + port

	fmt.Println("starting server...")

	mux := http.NewServeMux()
	apiConfig := &apiConfig{
		fileserverHits: atomic.Int32{},
	}
	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	// handlers
	mux.Handle("/app/", apiConfig.middlewareMetricInt(fileServerHandler))
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		numberOfHits := apiConfig.getHits()

		_, err := w.Write(fmt.Appendf(nil, "Hits: %d", numberOfHits))
		if err != nil {
			log.Fatal("error writing response body")
		}
	})
	mux.HandleFunc("POST /reset", func(w http.ResponseWriter, r *http.Request) {
		apiConfig.resetHits()
		w.WriteHeader(http.StatusOK)
	})

	// health check
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Fatal("error writing response body")
		}
	})

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
