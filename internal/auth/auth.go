package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(hash, password string) error {
	error := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return error
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chripy",
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if token.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// check if token has expired
		expiredTime, _ := token.Claims.GetExpirationTime()
		if time.Now().After(expiredTime.Local()) {
			return nil, fmt.Errorf("token has expired %v", expiredTime)
		}

		return []byte(tokenSecret), nil
	})

	if err != nil {
		fmt.Println(err)
		return uuid.Nil, err
	}

	id, err_claims := token.Claims.GetSubject()
	if err_claims != nil {
		fmt.Println(err_claims)
	}
	uuid, _ := uuid.Parse(id)
	return uuid, err
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if len(strings.Split(bearer, " ")) == 0 {
		return "", errors.New("bearer")
	}

	b := strings.Split(bearer, " ")
	if cap(b) < 2 {
		return "", errors.New("can not find bearer token")
	}

	return b[1], nil
}

func MakeRefreshToken() (string, error) {
	byte32 := make([]byte, 32)
	rand.Read(byte32)
	encoded := hex.EncodeToString(byte32)

	return encoded, nil
}
