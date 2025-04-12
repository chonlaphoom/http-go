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

func (cfg *apiConfig) middlewareMetricInt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		fmt.Println("app hit!")
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
	mux.Handle("GET /app/", apiConfig.middlewareMetricInt(fileServerHandler))
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		numberOfHits := apiConfig.getHits()

		_, err := w.Write(fmt.Appendf(nil, `
		<html>
		  <body>
			    <h1>Welcome, Chirpy Admin</h1>
					    <p>Chirpy has been visited %d times!</p>
			</body>
		</html>
		`, numberOfHits))

		if err != nil {
			log.Fatal("error writing response body")
		}
	})
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		apiConfig.resetHits()
		w.WriteHeader(http.StatusOK)
	})

	// health check
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
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
