package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pingcap/log"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
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

type CustomClaims struct {
	Claims jwt.Claims
}

func (cc *CustomClaims) MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	cc.Claims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chripy",
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc.Claims)
	return token.SignedString(tokenSecret)
}

func (cc *CustomClaims) ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, cc.Claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if token.Method.Alg() == "HS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		fmt.Println(err)
	}

	id, err_claims := token.Claims.GetSubject()
	if err_claims != nil {
		fmt.Println(err_claims)
	}
	uuid, _ := uuid.Parse(id)
	return uuid, err
}
