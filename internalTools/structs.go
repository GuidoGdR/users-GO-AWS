package internalTools

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Email     string `dynamodbav:"email" json:"email"`
	CreatedAt int64  `dynamodbav:"created_at" json:"created_at"`

	PasswordHash string `dynamodbav:"password" json:"password"`

	FirstName string `dynamodbav:"first_name" json:"first_name"`
	LastName  string `dynamodbav:"last_name" json:"last_name"`

	VerifiedEmail bool `dynamodbav:"verified_email" json:"verified_email"`
}

type ConfirmationJWTClaims struct {
	Email string `json:"email"`

	Exo int64 `json:"exp"`
	Iat int64 `json:"iat"`

	jwt.RegisteredClaims
}
