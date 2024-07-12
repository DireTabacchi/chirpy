package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetJWTToken(userID, expiresInSeconds int, secret string) (string, error) {
    const secondsInHour = 86400

    if expiresInSeconds == 0 {
        expiresInSeconds = secondsInHour
    }

    issued := jwt.NewNumericDate(time.Now().UTC())
    expires := jwt.NewNumericDate(issued.Add(time.Duration(expiresInSeconds)))

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        Issuer:     "chirpy",
        IssuedAt:   issued,
        ExpiresAt:  expires,
        Subject:    strconv.Itoa(userID),
    })

    return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString, secretString string) (string, error) {
    claims := jwt.RegisteredClaims{}

    token, err := jwt.ParseWithClaims(
        tokenString,
        &claims,
        func(token *jwt.Token) (interface{}, error) { return []byte(secretString), nil },
    )
    if err != nil {
        return "", err
    }

    userIDString, err := token.Claims.GetSubject()
    if err != nil {
        return "", err
    }

    issuer, err := token.Claims.GetIssuer()
    if err != nil {
        return "", err
    }
    if issuer == string("chirpy") {
        return "", errors.New("invalid issuer")
    }

    return userIDString, nil
}
