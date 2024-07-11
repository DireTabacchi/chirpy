package main

import (
    "log"
    "net/http"
    "encoding/json"
)

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

