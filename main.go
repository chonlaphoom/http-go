package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		decode_error := decoder.Decode(&params)

		// handle other error
		if decode_error != nil {
			responseWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		// handle error body exceed 140 chars
		maxBodySize := 140
		if len(params.Body) > maxBodySize {
			responseWithError(w, http.StatusBadRequest, "Chirp is too long")
			return
		}

		// handle success
		cleanInput := stringReplaceAll(params.Body, []string{"kerfuffle", "sharbert", "fornax"})
		type cleanRespond struct {
			Cleaned_body string `json:"cleaned_body"`
		}
		respondWithJSON(w, http.StatusOK, cleanRespond{Cleaned_body: cleanInput})
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

func jsonToByteWithMarshal[T any](js T) ([]byte, error) {
	data, err := json.Marshal(js)
	if err != nil {
		log.Printf("error marShalling json %s", err)
	}
	return data, err
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	body, errJsonToByte := jsonToByteWithMarshal(payload)
	if errJsonToByte != nil {
		log.Fatal("error encoding payload response")
	}

	_, err := w.Write(body)
	if err != nil {
		log.Fatal("error writing response body")
	}
}

func responseWithError(w http.ResponseWriter, code int, msg string) {
	type errorRes struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorRes{Error: msg})
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
