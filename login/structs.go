package main

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type smallUser struct {
	Email string `dynamodbav:"email" json:"email"`

	PasswordHash string `dynamodbav:"password" json:"password"`

	FirstName string `dynamodbav:"first_name" json:"first_name"`
	LastName  string `dynamodbav:"last_name" json:"last_name"`

	EmailVerified bool `dynamodbav:"email_verified" json:"email_verified"`
}
