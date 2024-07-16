package auth

import (
	"crypto/rand"
    "encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func HashPassword(pw string) (string, error) {
    passHash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }

    return string(passHash), nil
}

func CheckHashedPassword(password, hashedPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func MakeJWT(userID int, secret string, expiresIn time.Duration) (string, error) {
    signingKey := []byte(secret)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        Issuer:     "chirpy",
        IssuedAt:   jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt:  jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
        Subject:    strconv.Itoa(userID),
    })

    return token.SignedString(signingKey)
}

func MakeRefreshToken() (string, error) {
    tokenBytes := make([]byte, 32)

    _, err := rand.Read(tokenBytes)
    if err != nil {
        return "", err
    }

    return hex.EncodeToString(tokenBytes), nil
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
    if issuer != string("chirpy") {
        return "", errors.New("invalid issuer")
    }

    return userIDString, nil
}

func GetBearerToken(headers http.Header) (string, error) {
    authHeader := headers.Get("Authorization")
    if authHeader == "" {
        return "", ErrNoAuthHeaderIncluded
    }

    splitAuth := strings.Split(authHeader, " ")
    if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
        return "", errors.New("malformed authorization header")
    }

    return splitAuth[1], nil
}
