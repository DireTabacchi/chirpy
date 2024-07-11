package main

import (
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
    chirpIDStr := r.PathValue("id")
    chirpID, err := strconv.Atoi(chirpIDStr)
    if err != nil {
        log.Printf("Error converting to int: %v\n", err)
        respondWithError(w, http.StatusInternalServerError, "Testing some code")
        return
    }

    chirp, err := cfg.db.GetChirp(chirpID)

    if err != nil {
        respondWithError(w, http.StatusNotFound, "Chirp not found")
    }

    respondWithJSON(w, http.StatusOK, chirp)
}
