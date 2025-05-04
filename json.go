package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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

func jsonToByteWithMarshal[T any](js T) ([]byte, error) {
	data, err := json.Marshal(js)
	if err != nil {
		log.Printf("error marShalling json %s", err)
	}
	return data, err
}
