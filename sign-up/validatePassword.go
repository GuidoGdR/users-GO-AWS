package main

import (
	"fmt"
	"regexp"
)

//
//  Username     string `json:"username"`
//  PasswordHash string `json:"password"`
//  Email        string `json:"email"`
//  FirstName string `json:"first_name"`
//  LastName  string `json:"last_name"`

func validatePassword(password string) ([3]string, error) {

	errorTexts := [3]string{"", "", ""}
	length := len(password)

	if length < 8 || length > 32 {

		errorTexts[0] = "Contraseña muy larga o muy corta."
		errorTexts[1] = "Longitud: minima 8, maxima 32."
		errorTexts[2] = "password"
		err := fmt.Errorf("%s %s", errorTexts[0], errorTexts[1])

		return errorTexts, err
	}

	if !regexp.MustCompile(`[0-9]`).MatchString(password) {

		errorTexts[0] = "Contraseña insegura."
		errorTexts[1] = "La contraseña debe contener minimo un numero."
		errorTexts[2] = "password"
		err := fmt.Errorf("%s %s", errorTexts[0], errorTexts[1])

		return errorTexts, err
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(password) {

		errorTexts[0] = "Contraseña insegura."
		errorTexts[1] = "La contraseña debe contener minimo una letra minuscula."
		errorTexts[2] = "password"
		err := fmt.Errorf("%s %s", errorTexts[0], errorTexts[1])

		return errorTexts, err
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {

		errorTexts[0] = "Contraseña insegura."
		errorTexts[1] = "La contraseña debe contener minimo una letra mayuscula."
		errorTexts[2] = "password"
		err := fmt.Errorf("%s %s", errorTexts[0], errorTexts[1])

		return errorTexts, err
	}

	//if !regexp.MustCompile(`[!@#\$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password) {
	//	return errors.New("min one special character")
	//}

	return errorTexts, nil
}
