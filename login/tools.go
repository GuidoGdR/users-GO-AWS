package main

import (
	// base

	"time"

	// from third parties

	"github.com/golang-jwt/jwt/v5"

	"internalTools"

	"github.com/GuidoGdR/users-GO-AWS/usersTools"
)

func generateJWT(user *internalTools.User) (string, error) {
	claims := usersTools.JwtClaims{
		Email:     user.Email,
		CreatedAt: user.CreatedAt,

		FirstName: user.FirstName,
		LastName:  user.LastName,

		Exp: time.Now().Add(time.Hour * 24).Unix(),
		Iat: time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
