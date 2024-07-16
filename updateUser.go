package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/DireTabacchi/chirpy/internal/auth"
)

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Password    string `json:"password"`
        Email       string `json:"email"`
    }
    type response struct {
        User
    }

    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
        return
    }

    subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
    if err != nil {
        log.Printf("subject error: %v\n", err)
        respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
        return
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body")
        return
    }

    hashedPassword, err := auth.HashPassword(params.Password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
        return
    }

    userID, err := strconv.Atoi(subject)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
        return
    }

    user, err := cfg.db.UpdateUser(userID, params.Email, hashedPassword)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
        return
    }

    respondWithJSON(w, http.StatusOK, response{
        User{
            ID: user.ID,
            Email: user.Email,
        },
    })
}
