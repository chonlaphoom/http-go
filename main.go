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
	"time"

	"github.com/chonlaphoom/http-go/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             database.Queries
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
	// load env
	err_load_env := godotenv.Load()
	if err_load_env != nil {
		log.Fatal("error load env")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error open postgres dbUTRL:", dbURL)
	}

	apiConfig := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             *database.New(db),
	}

	mux := http.NewServeMux()
	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	// handlers
	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		type paramsT struct {
			Email string `json:"email"`
		}

		decoder := json.NewDecoder(r.Body)
		params := paramsT{}
		decode_error := decoder.Decode(&params)

		// handle decode error
		if decode_error != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Something went wrong during decode request body")
			return
		}
		fmt.Println("email", params.Email)
		_user, err := apiConfig.db.CreateUser(r.Context(), sql.NullString{Valid: params.Email != "", String: params.Email})
		if err != nil {
			fmt.Println("error creating user")
			fmt.Println(err)
		}

		user := &User{
			ID:        _user.ID,
			UpdatedAt: _user.UpdatedAt.Time,
			CreatedAt: _user.CreatedAt.Time,
			Email:     _user.Email.String,
		}

		RespondWithJSON(w, http.StatusCreated, user)
	})
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		platform := os.Getenv("PLATFORM")
		if platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "text/plain")

		err := apiConfig.db.ResetUsers(r.Context())
		if err != nil {
			fmt.Println("error during reset users table")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
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
	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		decode_error := decoder.Decode(&params)

		// handle other error
		if decode_error != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		// handle error body exceed 140 chars
		maxBodySize := 140
		if len(params.Body) > maxBodySize {
			ResponseWithError(w, http.StatusBadRequest, "Chirp is too long")
			return
		}

		// handle success
		cleanInput := stringReplaceAll(params.Body, []string{"kerfuffle", "sharbert", "fornax"})
		type cleanRespond struct {
			Cleaned_body string `json:"cleaned_body"`
		}
		RespondWithJSON(w, http.StatusOK, cleanRespond{Cleaned_body: cleanInput})
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

func jsonToByteWithMarshal[T any](js T) ([]byte, error) {
	data, err := json.Marshal(js)
	if err != nil {
		log.Printf("error marShalling json %s", err)
	}
	return data, err
}

// TODO write test
func stringReplaceAll(target string, profaneList []string) string {
	splitBy := " "
	lower := strings.ToLower(target)
	lower_splitted := strings.Split(lower, splitBy)
	target_splitted := strings.Split(target, splitBy)

	var indexToReplace []int
	for _, profane := range profaneList {
		for index, eachSplitted := range lower_splitted {
			if eachSplitted == profane {
				indexToReplace = append(indexToReplace, index)
			}
		}
	}

	for _, index := range indexToReplace {
		target_splitted[index] = "****"
	}

	return strings.Join(target_splitted, " ")
}
