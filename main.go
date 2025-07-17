package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/chonlaphoom/http-go/internal/database"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type UserWToken struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Token         string    `json:"token"`
	Refresh_token string    `json:"refresh_token"`
}

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             database.Queries
	tokenString    string
}

func main() {
	// load env
	dbURL, tokenStr := loadEnvironment()
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error open postgres dbUTRL:", dbURL)
	}

	apiConfig := &ApiConfig{
		FileserverHits: atomic.Int32{},
		Db:             *database.New(db),
		tokenString:    tokenStr,
	}

	mux := http.NewServeMux()
	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	// handlers
	mux.HandleFunc("POST /api/users", apiConfig.createUser)
	mux.HandleFunc("PUT /api/users", apiConfig.updateUser)
	mux.HandleFunc("POST /api/login", apiConfig.login)
	mux.HandleFunc("POST /api/refresh", apiConfig.refresh)
	mux.HandleFunc("POST /api/revoke", apiConfig.revoke)
	mux.HandleFunc("POST /admin/reset", apiConfig.resetUsers)
	mux.HandleFunc("GET /admin/metrics", apiConfig.metric)
	mux.HandleFunc("POST /api/chirps", apiConfig.createChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.getChirps)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiConfig.deleteChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.getChirps)
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.Handle("GET /app/", apiConfig.middlewareMetricInt(fileServerHandler))

	port := "8080"
	address := ":" + port

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	fmt.Println("starting server...")

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
