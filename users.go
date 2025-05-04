package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func (cfg *ApiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type paramsT struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paramsT{}
	decode_error := decoder.Decode(&params)

	// handle decode error
	if decode_error != nil {
		responseWithError(w, http.StatusInternalServerError, "Something went wrong during decode request body")
		return
	}
	fmt.Println("email", params.Email)
	_user, err := cfg.Db.CreateUser(r.Context(), sql.NullString{Valid: params.Email != "", String: params.Email})
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

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *ApiConfig) resetUsers(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	err := cfg.Db.ResetUsers(r.Context())
	if err != nil {
		fmt.Println("error during reset users table")
	}
	w.WriteHeader(http.StatusOK)

	_, error_writing := w.Write([]byte("OK"))

	if error_writing != nil {
		fmt.Println("error writing response")
	}
}
