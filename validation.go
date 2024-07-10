package main

import (
	"encoding/json"
    "log"
	"net/http"
    "strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
    type chirpPost struct {
        Body string `json:"body"`
    }

    type validResponse struct {
        CleanedBody string `json:"cleaned_body"`
    }

    decoder := json.NewDecoder(r.Body)
    chirp := chirpPost{}
    err := decoder.Decode(&chirp)
    if err != nil {
        log.Printf("Error decoding chirp: %s", err)
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode chirp")
        return
    }

    if len(chirp.Body) > 140 {
        log.Printf("Chirp too long")
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }

    profanities := map[string]struct{}{
        "kerfuffle": {},
        "sharbert": {},
        "fornax": {},
    }

    body := getCleanBody(chirp.Body, profanities) 

    respondWithJSON(w, http.StatusOK, validResponse{
        CleanedBody: body,
    })
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
    if code > 499 {
        log.Printf("Responding with 5XX error: %s", msg)
    }

    type errResponse struct {
        Error string `json:"error"`
    }
    respondWithJSON(w, code, errResponse{
        Error: msg,
    })
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    data, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marhalling JSON: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(data)
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
