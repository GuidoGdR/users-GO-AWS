package usersTools

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Exp int64 `json:"exp"` // expiration time
	Iat int64 `json:"iat"` // emition time

	jwt.RegisteredClaims
}

func DecodeJWT(token string, jwtKey string) (JwtClaims, error) {

	// Parse token
	tkn, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return JwtClaims{}, err
	}

	claims, ok := tkn.Claims.(JwtClaims)

	if ok && tkn.Valid {

		return claims, nil
	}

	return JwtClaims{}, errors.New("invalid token")

}
