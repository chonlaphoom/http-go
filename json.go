package main

import "log"
import (
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
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

func ResponseWithError(w http.ResponseWriter, code int, msg string) {
	type errorRes struct {
		Error string `json:"error"`
	}
	RespondWithJSON(w, code, errorRes{Error: msg})
}
