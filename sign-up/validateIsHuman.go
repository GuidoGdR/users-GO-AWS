package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type turnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
	Action      *string  `json:"action,omitempty"`
	Cdata       *string  `json:"cdata,omitempty"`
}

func validateIsHuman(token string) (bool, string, error) {

	data := map[string]string{
		"secret":   turnstileKey,
		"response": token,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, "Error al generar la consulta para validar si se es humano", err
	}

	resp, err := http.Post("https://challenges.cloudflare.com/turnstile/v0/siteverify", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, "Error al intentar consultar externamente si se es humano", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "Error al intentar leer la respuesta de si se es humano", err
	}

	var turnstileResp turnstileResponse
	err = json.Unmarshal(body, &turnstileResp)
	if err != nil {
		return false, "Error al intentar deserializar el json de la respuesta de si se es humano", err
	}

	return turnstileResp.Success, "", nil
}
