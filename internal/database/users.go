package database

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID int          `json:"id"`
    Email string    `json:"email"`
    Password []byte `json:"password"`
}

type UserResponse struct {
    ID int          `json:"id"`
    Email string    `json:"email"`
}

func (db *DB) CreateUser(email string, password string) (UserResponse, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return UserResponse{}, err
    }

    id := len(dbStructure.Users) + 1

    passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return UserResponse{}, err
    }

    user := User{ ID: id, Email: email, Password: passHash }
    dbStructure.Users[id] = user

    err = db.writeDB(dbStructure)
    if err != nil {
        return UserResponse{}, err
    }

    respUser := UserResponse{ ID: id, Email: email }
    return respUser, nil
}

func (db *DB) GetUser(email string, password string) (UserResponse, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return UserResponse{}, err
    }

    user := User{}
    for _, u := range dbStructure.Users {
        if u.Email == email {
            user = u
        }
    }

    if user.Email == "" && user.ID == 0 {
        return UserResponse{}, ErrNotExist
    }

    err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
    if err != nil {
        return UserResponse{}, err
    }

    respUser := UserResponse { ID: user.ID, Email: user.Email }

    return respUser, nil
}
