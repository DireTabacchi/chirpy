package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) postUserHandler(w http.ResponseWriter, r *http.Request) {
    type userPost struct {
        Email string `json:"email"`
    }

    decoder := json.NewDecoder(r.Body)
    user := userPost{}

    err := decoder.Decode(&user)
    if err != nil {
        log.Printf("Error decoding user: %s\n", err)
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode user")
        return
    }

    postedUser, err := cfg.db.CreateUser(user.Email)
    if err != nil {
        log.Printf("Error creating user: %v\n", err)
        respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
        return
    }

    respondWithJSON(w, http.StatusCreated, postedUser)
}
