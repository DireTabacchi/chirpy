package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
    type userPost struct {
        Email string `json:"email"`
        Password string `json:"password"`
    }

    decoder := json.NewDecoder(r.Body)
    user := userPost{}

    err := decoder.Decode(&user)
    if err != nil {
        log.Printf("Error decoding user login: %v\n", err)
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode login")
        return
    }

    result, err := cfg.db.GetUser(user.Email, user.Password)
    if err == bcrypt.ErrMismatchedHashAndPassword {
        respondWithError(w, http.StatusUnauthorized, "Incorrect password")
        return
    } else if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, result)
}
