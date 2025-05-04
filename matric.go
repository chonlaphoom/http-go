package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *ApiConfig) metric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	numberOfHits := cfg.getHits()

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
}

func (cfg *ApiConfig) getHits() int32 {
	return cfg.FileserverHits.Load()
}

func (cfg *ApiConfig) resetHits() {
	cfg.FileserverHits.Store(0)
}

func (cfg *ApiConfig) middlewareMetricInt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		fmt.Println("app hit!")
		next.ServeHTTP(w, r)
	})
}
