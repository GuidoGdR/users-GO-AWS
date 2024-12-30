package internalTools

import (
	"encoding/json"
)

func MakeErrorBody(title string, message string, field string) (string, error) {

	body, err := json.Marshal(map[string]string{
		"error":   title,
		"message": message,
		"field":   field,
	})

	if err != nil {
		return "", err
	}

	return string(body), nil
}
