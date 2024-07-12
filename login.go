package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
    "github.com/DireTabacchi/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
    type userPost struct {
        Email string            `json:"email"`
        Password string         `json:"password"`
        ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
    }

    type userResponse struct {
        Email string    `json:"email"`
        ID int          `json:"id"`
        Token string    `json:"token"`
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

    signedJWT, err := auth.GetJWTToken(result.ID, user.ExpiresInSeconds, cfg.jwtSecret)
    if err != nil {
        log.Printf("Error getting signed token: %v/n", err)
        respondWithError(w, http.StatusInternalServerError, "Couldn't get signed token")
        return
    }

    resp := userResponse{
        Email: result.Email,
        ID: result.ID,
        Token: signedJWT,
    }

    log.Printf("Got token: %s\n", signedJWT)

    respondWithJSON(w, http.StatusOK, resp)
}
