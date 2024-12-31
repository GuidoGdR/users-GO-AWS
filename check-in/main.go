package main

// Required env vars:
// USERS_TABLE
// JWT_SECRET_2

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

	"github.com/golang-jwt/jwt/v5"

	// custom
	"internalTools"
)

var dbSvc *dynamodb.Client
var usersTableName string

var jwtKey []byte

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Logs
	log.SetPrefix("check-in ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// validate request
	var req Request

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {

		body, _ := internalTools.MakeErrorBody("Cuerpo de solicitud inválido", "Requiere token", "token")
		return internalTools.Response(400, body), nil
	}

	if req.Token == "" {

		body, _ := internalTools.MakeErrorBody("Token invalido", "no se pudo validar el token", "token")
		return internalTools.Response(400, body), nil
	}

	// Parse token
	tkn, err := jwt.ParseWithClaims(req.Token, &internalTools.ConfirmationJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		body, _ := internalTools.MakeErrorBody("Token invalido", "no se pudo validar el token", "token")
		return internalTools.Response(400, body), nil
	}

	claims, ok := tkn.Claims.(*internalTools.ConfirmationJWTClaims)

	if !ok || !tkn.Valid {

		body, _ := internalTools.MakeErrorBody("Token invalido", "no se pudo validar el token", "token")
		return internalTools.Response(400, body), nil
	}

	// update "verifiedEmail"
	if userDontExist, errMsg, err := updateVerifiedEmail(ctx, claims.Email); userDontExist {

		body, _ := internalTools.MakeErrorBody("Token invalido", "no se pudo validar el token", "token")
		return internalTools.Response(400, body), err

	} else if err != nil {

		log.Print(fmt.Sprintf("%s\n", errMsg), err)
		body, _ := internalTools.MakeErrorBody("Error interno", errMsg, "")
		return internalTools.Response(500, body), err
	}

	jsonResp, _ := json.Marshal(map[string]string{"message": fmt.Sprintf("Email %s confirmado correctamente", claims.Email)})

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

	jwtKey = []byte(os.Getenv("JWT_SECRET_2"))
	if len(jwtKey) == 0 {
		panic("JWT_SECRET_2 environment variable not set")
	}

	// AWS
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	// DB
	dbSvc = dynamodb.NewFromConfig(cfg)
}
