package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chonlaphoom/http-go/internal/auth"
	"github.com/chonlaphoom/http-go/internal/database"
	"github.com/google/uuid"
)

type chirpParam struct {
	Body string `json:"body"`
}

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	User_id   uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) createChirps(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	bodyParams := chirpParam{}
	decode_error := decoder.Decode(&bodyParams)

	// handle other error
	if decode_error != nil {
		responseWithError(w, http.StatusInternalServerError, "Something went wrong, decoding body")
		return
	}

	// unauthorized
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "can not get bearer token")
		return
	}

	userId, errValidate := auth.ValidateJWT(bearer, cfg.tokenString)

	if errValidate != nil {
		responseWithError(w, http.StatusUnauthorized, "token is not valid")
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

	newChirp, error_create_chirp := cfg.Db.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanBody, UserID: userId})

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

func (cfg *ApiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")

	_chirps, err := cfg.Db.GetChirps(r.Context())

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "something went wrong getting chirps")
		return
	}

	chirps := []Chirp{}
	for _, chirp := range _chirps {
		chirps = append(chirps, Chirp{Body: chirp.Body,
			User_id:   chirp.UserID,
			CreatedAt: chirp.CreatedAt.Time,
			UpdatedAt: chirp.UpdatedAt.Time,
			Id:        chirp.ID,
		})
	}

	if chirpId != "" {
		found := Chirp{}
		isFound := false
		for _, v := range chirps {
			if string(v.Id.String()) == chirpId {
				isFound = true
				found = v
				break
			}
		}
		if isFound {
			respondWithJSON(w, http.StatusOK, found)
			return
		}
		responseWithError(w, http.StatusNotFound, "not found")
	} else {
		// get all chirps
		respondWithJSON(w, http.StatusOK, chirps)
	}
}
