package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/DireTabacchi/chirpy/internal/auth"
	"github.com/DireTabacchi/chirpy/internal/database"
)

type User struct {
    ID          int     `json:"id"`
    Email       string  `json:"email"`
    Password    string  `json:"-"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Email       string `json:"email"`
        Password    string `json:"password"`
    }
    type response struct {
        User
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}

    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding user: %s\n", err)
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode user")
        return
    }

    hashedPass, err := auth.HashPassword(params.Password)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
        return
    }

    user, err := cfg.db.CreateUser(params.Email, hashedPass)
    if err != nil {
        if errors.Is(err, database.ErrAlreadyExists) {
            respondWithError(w, http.StatusConflict, "User already exists")
            return
        }

        respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
        return
    }


    respondWithJSON(w, http.StatusCreated, response{
        User{
            ID: user.ID,
            Email: user.Email,
        },
    })
}
