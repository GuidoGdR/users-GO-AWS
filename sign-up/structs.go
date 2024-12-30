package main

type Request struct {
	Email   string `json:"email"`
	IsHuman string `json:"is_human"`

	Password string `json:"password"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
