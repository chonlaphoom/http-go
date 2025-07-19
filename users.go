package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/chonlaphoom/http-go/internal/auth"
	"github.com/chonlaphoom/http-go/internal/database"
	"github.com/google/uuid"
)

type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *ApiConfig) login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := UserData{}
	decode_error := decoder.Decode(&params)

	// handle decode error
	if decode_error != nil {
		responseWithError(w, http.StatusInternalServerError, "Something went wrong during decode request body")
		return
	}

	dbUser, err := cfg.Db.GetUserByEmail(r.Context(), sql.NullString{Valid: params.Email != "", String: params.Email})

	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Something went wrong getting user")
		return
	}

	if err := auth.CheckPasswordHash(dbUser.HashedPassword, params.Password); err != nil {
		responseWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	expires := 60 * 60
	token, errCreateToken := auth.MakeJWT(dbUser.ID, cfg.tokenString, time.Second*time.Duration(expires))
	refresh_token, errCreateRefreshToken := auth.MakeRefreshToken()

	if errCreateToken != nil {
		fmt.Println(errCreateToken)
		responseWithError(w, http.StatusUnauthorized, "something went wrong during create token")
		return
	}
	if errCreateRefreshToken != nil {
		fmt.Println(errCreateRefreshToken)
		responseWithError(w, http.StatusUnauthorized, "something went wrong during create refresh token")
		return
	}

	t, errInsertRefreshToken := cfg.Db.InsertRefreshToken(r.Context(), database.InsertRefreshTokenParams{UserID: dbUser.ID, Token: refresh_token})

	if errInsertRefreshToken != nil {
		fmt.Println(errInsertRefreshToken)
		responseWithError(w, http.StatusUnauthorized, "something went wrong during insert refresh token")
	}

	user := UserWToken{
		ID:            dbUser.ID,
		UpdatedAt:     dbUser.UpdatedAt.Time,
		CreatedAt:     dbUser.CreatedAt.Time,
		Email:         dbUser.Email.String,
		Token:         token,
		Refresh_token: t.Token,
		IsChirpyRed:   dbUser.IsChirpyRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *ApiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("begin create user...")

	decoder := json.NewDecoder(r.Body)
	params := UserData{}
	decode_error := decoder.Decode(&params)

	// handle decode error
	if decode_error != nil {
		fmt.Printf("err: %e", decode_error)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong during decode request body")
		return
	}

	pass, err := auth.HashPassword(params.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Something went wrong hashing password")
		return
	}

	createUserParams := database.CreateUserParams{
		Email:          sql.NullString{Valid: params.Email != "", String: params.Email},
		HashedPassword: pass,
	}
	_user, err := cfg.Db.CreateUser(r.Context(), createUserParams)

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

func (cfg *ApiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetAPIKey(r.Header)
	if err != nil {
		fmt.Print(bearer)
		responseWithError(w, http.StatusUnauthorized, "can not get refresh token")
		return
	}

	refToken, errGettingRefreshToken := cfg.Db.GetRefreshTokenByToken(r.Context(), bearer)
	if errGettingRefreshToken != nil {
		responseWithError(w, http.StatusUnauthorized, "error refresh token not found")
		return
	}

	isExpire := time.Now().After(refToken.ExpiresAt.Time)
	if isExpire {
		responseWithError(w, http.StatusUnauthorized, "error refresh token expired")
		return
	}

	if refToken.RevokedAt.Valid {
		responseWithError(w, http.StatusUnauthorized, "error refresh token already revoked")
		return
	}

	expires := 60 * 60
	access_token, errCreateToken := auth.MakeJWT(refToken.UserID, cfg.tokenString, time.Second*time.Duration(expires))
	if errCreateToken != nil {
		responseWithError(w, http.StatusUnauthorized, "error create acess refresh token")
	}

	type res struct {
		Token string `json:"token"`
	}
	response := res{Token: access_token}
	respondWithJSON(w, http.StatusOK, &response)
}

func (cfg *ApiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetAPIKey(r.Header)

	if err != nil {
		fmt.Print(bearer)
		responseWithError(w, http.StatusUnauthorized, "can not revoke refresh token")
		return
	}

	rt, erroGetRT := cfg.Db.GetRefreshTokenByToken(r.Context(), bearer)

	if erroGetRT != nil {
		responseWithError(w, http.StatusUnauthorized, "can not get refresh token")
		return
	}
	if rt.RevokedAt.Valid {
		// already revoke
		responseWithError(w, http.StatusUnauthorized, "already revoked")
		return
	}

	erroRevoke := cfg.Db.RevokeExistingRefreshToken(r.Context(), bearer)
	if erroRevoke != nil {
		responseWithError(w, http.StatusUnauthorized, "can not get refresh token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}

func (cfg *ApiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetAPIKey(r.Header)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusUnauthorized, "can not get access token")
		return
	}

	inputParam := UserData{}
	newDecoder := json.NewDecoder(r.Body)
	decoderError := newDecoder.Decode(&inputParam)
	if decoderError != nil {
		fmt.Println(decoderError)
		responseWithError(w, http.StatusUnauthorized, "can decode input from client")
		return
	}

	userId, errValidate := auth.ValidateJWT(bearer, cfg.tokenString)
	if errValidate != nil {
		responseWithError(w, http.StatusUnauthorized, "error validating bearer token")
		return
	}

	hashedPass, errorHash := auth.HashPassword(inputParam.Password)
	if errorHash != nil {
		responseWithError(w, http.StatusUnauthorized, "error hashing password")
		return
	}

	updateData := database.UpdateUserByUserIdParams{
		Email:          sql.NullString{Valid: true, String: inputParam.Email},
		HashedPassword: hashedPass,
		ID:             userId,
	}
	user, errorUpdateUser := cfg.Db.UpdateUserByUserId(r.Context(), updateData)
	if errorUpdateUser != nil {
		fmt.Println(errorUpdateUser)
		responseWithError(w, http.StatusUnauthorized, "error updating user database")
		return
	}

	type ResponseUser struct {
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	respondWithJSON(w, http.StatusOK, ResponseUser{
		Id:        userId.String(),
		Email:     user.Email.String,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	})
}

func (cfg *ApiConfig) updateChirpRed(w http.ResponseWriter, r *http.Request) {
	apiKey, errGettingKey := auth.GetAPIKey(r.Header)

	if apiKey != cfg.polkaKey || errGettingKey != nil {
		fmt.Println("error: getting api keys polka", errGettingKey)
		respondWithJSON(w, http.StatusUnauthorized, nil)
		return
	}

	type updateChirpRed struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	body := updateChirpRed{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		fmt.Println("error: parsing webhook body", err)
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	if body.Event != "user.upgraded" {
		fmt.Println("error: event not user upgrade", body.Event)
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	userIdStr, errParseUserId := uuid.Parse(body.Data.UserId)
	if errParseUserId != nil {
		fmt.Println("error: parsing userId")
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	_, errUpdateUserChirpRed := cfg.Db.UpdateUserChirpRedTRUE(r.Context(), userIdStr)

	if errUpdateUserChirpRed != nil {
		fmt.Println("error: update user chirp red")
		respondWithJSON(w, http.StatusNotFound, nil)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
