package main

import (
	"encoding/json"
	"net/http"
	"slices"
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
	bearer, err := auth.GetAPIKey(r.Header)
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
		for _, c := range chirps {
			if string(c.Id.String()) == chirpId {
				isFound = true
				found = c
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
		author_id := r.URL.Query().Get("author_id")
		if author_id != "" {
			// filter by author id
			filterdChirps := []Chirp{}
			for _, chirp := range chirps {
				if chirp.User_id.String() == author_id {
					filterdChirps = append(filterdChirps, chirp)
				}
			}

			chirps = filterdChirps
		}

		descSort := r.URL.Query().Get("sort")
		if descSort == "desc" {
			// sort by desc
			slices.Reverse(chirps)
		}
		respondWithJSON(w, http.StatusOK, chirps)
	}
}

func (cfg *ApiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	token, errorGettingToken := auth.GetAPIKey(r.Header)
	chirpId := r.PathValue("chirpID")

	if errorGettingToken != nil {
		responseWithError(w, http.StatusUnauthorized, errorGettingToken.Error())
		return
	}

	userId, errorValidateJWT := auth.ValidateJWT(token, cfg.tokenString)
	if errorValidateJWT != nil {
		responseWithError(w, http.StatusForbidden, "Error validate JWT")
		return
	}

	id, _ := uuid.Parse(chirpId)
	chirp, errGetChirp := cfg.Db.GetChirpFromId(r.Context(), id)

	if errGetChirp != nil {
		responseWithError(w, http.StatusForbidden, "can not get chirp")
		return
	}

	if userId.String() != chirp.UserID.String() {
		responseWithError(w, http.StatusForbidden, "not your chirp!")
		return
	}

	err := cfg.Db.DeleteChirpsById(r.Context(), id)

	if err != nil {
		responseWithError(w, http.StatusNotFound, "something went wrong during delete chirp")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
