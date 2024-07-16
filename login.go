package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DireTabacchi/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Email       string `json:"email"`
        Password    string `json:"password"`
    }

    type response struct {
        User
        Token           string `json:"token"`
        RefreshToken    string `json:"refresh_token"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode login")
        return
    }

    user, err := cfg.db.GetUserByEmail(params.Email)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
    }

    err = auth.CheckHashedPassword(params.Password, user.HashedPassword)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Invalid password")
        return
    }

    signedJWT, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
    if err != nil {
        log.Printf("Error getting signed token: %v/n", err)
        respondWithError(w, http.StatusInternalServerError, "Couldn't get signed token")
        return
    }

    refreshToken, err := auth.MakeRefreshToken()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't make refresh token")
        return
    }

    err = cfg.db.SaveRefreshToken(user.ID, refreshToken)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token")
    }

    respondWithJSON(w, http.StatusOK, response{
        User{
            Email: user.Email,
            ID: user.ID,
        },
        signedJWT,
        refreshToken,
    })
}
