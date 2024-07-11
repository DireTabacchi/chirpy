package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) postChirpHandler(w http.ResponseWriter, r *http.Request) {
    type chirpPost struct {
        Body string `json:"body"`
    }

    decoder := json.NewDecoder(r.Body)
    chirp := chirpPost{}
    err := decoder.Decode(&chirp)
    if err != nil {
        log.Printf("Error decoding chirp: %s\n", err)
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode chirp")
        return
    }

    body, err := validateChirp(chirp.Body)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    postedChirp, err := cfg.db.CreateChirp(body)
    if err != nil {
        log.Printf("Error creating chirp: %v\n", err)
        respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
        return
    }

    respondWithJSON(w, http.StatusCreated, postedChirp)
}

func validateChirp(body string) (string, error) {
    const maxChirpLength = 140
    if len(body) > maxChirpLength {
        log.Printf("Chirp too long")
        return "", errors.New("Chirp is too long")
    }

    profanities := map[string]struct{}{
        "kerfuffle": {},
        "sharbert": {},
        "fornax": {},
    }

    cleaned := getCleanBody(body, profanities)
    return cleaned, nil
}

func getCleanBody(body string, profanities map[string]struct{}) string {
    words := strings.Split(body, " ")
    for i, word := range words {
        loweredWord := strings.ToLower(word)
        if _, ok := profanities[loweredWord]; ok {
            words[i] = "****"
        }
    }
    return strings.Join(words, " ")
}
