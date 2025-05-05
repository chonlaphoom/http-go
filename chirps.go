package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chonlaphoom/http-go/internal/database"
	"github.com/google/uuid"
)

type chirpParam struct {
	Body    string    `json:"body"`
	User_id uuid.UUID `json:"user_id"`
}

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	User_id   uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) chirps(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	bodyParams := chirpParam{}
	decode_error := decoder.Decode(&bodyParams)

	// handle other error
	if decode_error != nil {
		responseWithError(w, http.StatusInternalServerError, "Something went wrong, decoding body")
		return
	}

	// handle error body exceed 140 chars
	maxBodySize := 140
	if len(bodyParams.Body) > maxBodySize {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// handle success
	cleanBody := stringReplaceAll(bodyParams.Body, []string{"kerfuffle", "sharbert", "fornax"})

	newChirp, error_create_chirp := cfg.Db.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanBody, UserID: bodyParams.User_id})

	if error_create_chirp != nil {
		responseWithError(w, http.StatusBadRequest, "Can not create chirp, something went wrong!")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		Body:      newChirp.Body,
		User_id:   newChirp.UserID,
		CreatedAt: newChirp.CreatedAt.Time,
		UpdatedAt: newChirp.UpdatedAt.Time,
		Id:        newChirp.ID,
	})
}
