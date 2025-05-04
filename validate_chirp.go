package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validate_chirp(w http.ResponseWriter, r *http.Request) {
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
