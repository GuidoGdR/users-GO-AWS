package main

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

func validateEmail(email string) ([3]string, error) {

	errorTexts := [3]string{"Email invalido.", "El email no encaja con los formatos tradicionales.", "email"}
	returnable_error := fmt.Errorf("%s %s", errorTexts[0], errorTexts[1])
	email = strings.TrimSpace(email)

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return errorTexts, returnable_error
	}

	if !regexp.MustCompile(`^(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*)@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`).MatchString(email) {
		return errorTexts, returnable_error
	}

	if addr.Address != email {
		return errorTexts, returnable_error
	}

	return [3]string{"", "", ""}, nil
}
