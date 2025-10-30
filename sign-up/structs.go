package main

type Request struct {
	Email   string `json:"email"`
	IsHuman string `json:"is_human"`

	Password string `json:"password"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type smallUser struct {
	Email string `dynamodbav:"email" json:"email"`

	PasswordHash string `dynamodbav:"password" json:"password"`

	EmailVerified bool `dynamodbav:"verified_email" json:"verified_email"`
}
