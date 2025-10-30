package main

// Required env vars:
// JWT_SECRET
// USERS_TABLE

import (

	// base
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	// from third parties
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"golang.org/x/crypto/bcrypt"

	// custom
	"internalTools"
)

var dbSvc *dynamodb.Client
var usersTableName string

var jwtKey []byte

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Logs
	log.SetPrefix("login ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// validate request
	var req Request
	err := json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		body, _ := internalTools.MakeErrorBody("Cuerpo de solicitud inválido", "Requiere email y contraseña", "email y/o password")
		return internalTools.Response(400, body), nil
	}

	if req.Email == "" || req.Password == "" {

		body, _ := internalTools.MakeErrorBody("Credenciales incorrectas", "Email y/o contraseña invalidos", "email y/o password")
		return internalTools.Response(401, body), nil
	}

	user, err := getUser(ctx, req.Email)
	if err != nil {

		log.Print("Error al obtener el usuario.\n", err)
		body, _ := internalTools.MakeErrorBody("Error interno", "Error al obtener el usuario.", "")
		return internalTools.Response(500, body), err
	}

	if user == nil || !user.VerifiedEmail {

		body, _ := internalTools.MakeErrorBody("Credenciales incorrectas", "Email y/o contraseña invalidos", "email y/o password")
		return internalTools.Response(401, body), nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {

		body, _ := internalTools.MakeErrorBody("Credenciales incorrectas", "Email y/o contraseña invalidos", "email y/o password")
		return internalTools.Response(401, body), nil
	}

	token, err := generateJWT(user)
	if err != nil {

		log.Print("Error al generar el JWT.\n", err)
		body, _ := internalTools.MakeErrorBody("Error interno", "Error al generar el JWT", "")
		return internalTools.Response(500, body), err
	}

	jsonResp, _ := json.Marshal(map[string]string{
		"message": fmt.Sprintf("Bienvenido, %s!", user.Email),
		"token":   token,
	})

	return internalTools.Response(200, string(jsonResp)), nil
}

func main() {
	lambda.Start(HandleRequest)
}

func init() {
	// Environment Variables
	usersTableName = os.Getenv("USERS_TABLE")
	if usersTableName == "" {
		panic("USERS_TABLE environment variable not set")
	}

	jwtKey = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtKey) == 0 {
		panic("JWT_SECRET environment variable not set")
	}

	// AWS
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	// DB
	dbSvc = dynamodb.NewFromConfig(cfg)
}
