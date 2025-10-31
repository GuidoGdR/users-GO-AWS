package main

// Required env vars:
// USERS_TABLE
// SENDER_EMAIL_ADDRESS
// JWT_SECRET_2
// TURNSTILE_SECRET_KEY
// CONFIRMATION_URL

import (

	// base
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	// from third parties
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ses"

	"golang.org/x/crypto/bcrypt"

	// custom
	"internalTools"
)

var sesSvc *ses.Client
var senderEmailAddress string

var dbSvc *dynamodb.Client
var usersTableName string

var jwtKey []byte
var turnstileKey string

var confirmationUrl string

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// validate request
	var req Request

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {

		body, _ := internalTools.MakeErrorBody("Cuerpo de solicitud inválido", "Requiere token, email y contraseña", "is_human, email y/o password")
		return internalTools.Response(400, body), nil
	}

	// validate data
	if errorTexts, err := validateNewUser(&req); err != nil {

		body, _ := internalTools.MakeErrorBody(errorTexts[0], errorTexts[1], errorTexts[2])
		return internalTools.Response(400, body), nil

	}

	// validate is_human
	if turnstileKey != "" {
		if len(req.IsHuman) < 50 {

			body, _ := internalTools.MakeErrorBody("Token invalido", "No se pudo validar que no sea un robot.", "is_human")
			return internalTools.Response(400, body), nil
		}

		isHuman, errMsg, err := validateIsHuman(req.IsHuman)

		if err != nil {

			log.Print(fmt.Sprintf("%s\n", errMsg), err)
			body, _ := internalTools.MakeErrorBody("Error interno", errMsg, "")
			return internalTools.Response(500, body), err
		}

		if !isHuman {

			body, _ := internalTools.MakeErrorBody("Token invalido", "No se pudo validar que no sea un robot.", "is_human")
			return internalTools.Response(400, body), nil
		}
	}

	// make password hash

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {

		log.Print("Error al encriptar la contraseña.\n", err)
		body, _ := internalTools.MakeErrorBody("Error interno", "Error al encriptar la contraseña", "")
		return internalTools.Response(500, body), err
	}

	// user obj
	newUser := internalTools.User{
		Email:     req.Email,
		CreatedAt: time.Now().Unix(),

		PasswordHash: string(passwordHash),

		FirstName: req.FirstName,
		LastName:  req.LastName,

		VerifiedEmail: false,
	}

	// check duplicates
	smallUser, err := getUser(ctx, newUser.Email)

	if err != nil {

		log.Print("Error al comprobar el uso previo del email.\n", err)
		body, _ := internalTools.MakeErrorBody("Error interno", "Error al comprobar el uso previo del email.", "")
		return internalTools.Response(500, body), err

	}

	if smallUser != nil {

		// have a previous User
		if smallUser.EmailVerified {

			body, _ := internalTools.MakeErrorBody("Email no disponible.", "Ya existe una cuenta vinculada a ese email.", "email")
			return internalTools.Response(400, body), nil
		}
		// previous User not verified

		if smallUser.PasswordHash != newUser.PasswordHash {

			// uptade pasword
			if errorText, err := updatePassword(ctx, newUser.Email, newUser.PasswordHash); err != nil {

				log.Print(fmt.Sprintf("%s\n", errorText), err)
				body, _ := internalTools.MakeErrorBody("Error interno", errorText, "")
				return internalTools.Response(500, body), err
			}
			smallUser.PasswordHash = newUser.PasswordHash
		}

	} else {
		// save user
		if errorText, err := createUser(ctx, newUser); err != nil {

			log.Print(fmt.Sprintf("%s\n", errorText), err)
			body, _ := internalTools.MakeErrorBody("Error interno", errorText, "")
			return internalTools.Response(500, body), err
		}
	}

	// send mail
	jwt, err := generateJWT(newUser.Email)
	if err != nil {

		log.Print("Error al intentar generar el token para confirmación de email.\n", err)
		body, _ := internalTools.MakeErrorBody("Error interno", "Error al intentar generar el token para confirmación de email.", "")
		return internalTools.Response(500, body), err
	}

	if err = sendConfirmationToEmail(ctx, newUser, fmt.Sprintf("%s?token=%s", confirmationUrl, jwt)); err != nil {

		log.Print("Error al intentar enviar el email de confirmación.\n", err)
		body, _ := internalTools.MakeErrorBody("Error interno", "Error al intentar enviar el email de confirmación.", "")
		return internalTools.Response(500, body), err
	}

	jsonResp, _ := json.Marshal(map[string]string{"message": fmt.Sprintf("Se a enviado un email de confirmación a %s.", req.Email)})

	return internalTools.Response(200, string(jsonResp)), nil
}

func main() {
	lambda.Start(HandleRequest)
}

func init() {
	// Logs
	log.SetPrefix("sign-up ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Environment Variables

	usersTableName = os.Getenv("USERS_TABLE")
	if usersTableName == "" {
		panic("USERS_TABLE environment variable not set")
	}

	senderEmailAddress = os.Getenv("SENDER_EMAIL_ADDRESS")
	if senderEmailAddress == "" {
		panic("SENDER_EMAIL_ADDRESS environment variable not set")
	}

	jwtKey = []byte(os.Getenv("JWT_SECRET_2"))
	if len(jwtKey) == 0 {
		panic("JWT_SECRET_2 environment variable not set")
	}

	turnstileKey = os.Getenv("TURNSTILE_SECRET_KEY")
	if turnstileKey == "" {
		log.Print("Turnstile disabled\n")
	}

	confirmationUrl = os.Getenv("CONFIRMATION_URL")
	if confirmationUrl == "" {
		panic("CONFIRMATION_URL environment variable not set")
	}

	//
	// AWS
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	// DB
	dbSvc = dynamodb.NewFromConfig(cfg)

	// SES
	sesSvc = ses.NewFromConfig(cfg)
}
