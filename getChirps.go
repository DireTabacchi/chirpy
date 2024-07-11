package main

import (
    "log"
	"net/http"
    "sort"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
    chirps, err := cfg.db.GetChirps()

    if err != nil {
        log.Printf("Error while getting chirps from database: %v\n", err)
        respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
        return
    }

    sort.Slice(chirps, func(i, j int) bool {
        return chirps[i].ID < chirps[j].ID
    })

    respondWithJSON(w, http.StatusOK, chirps)
}
