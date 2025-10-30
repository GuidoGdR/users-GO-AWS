package main

import (
	// base
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/ses/types"

	"github.com/golang-jwt/jwt/v5"

	"internalTools"
)

func validateNewUser(request *Request) ([3]string, error) {

	if errorTexts, err := validateEmail(request.Email); err != nil {
		return errorTexts, err
	}

	if errorTexts, err := validatePassword(request.Password); err != nil {
		return errorTexts, err
	}

	if len(request.FirstName) > 48 {
		errorTexts := [3]string{
			"Nombre demasiado largo.",
			"Longitud: maxima 48.",
			"first_name",
		}
		err := fmt.Errorf("%s %s", errorTexts[0], errorTexts[1])
		return errorTexts, err
	}

	if len(request.LastName) > 48 {
		errorTexts := [3]string{
			"Apellido demasiado largo.",
			"Longitud: maxima 48.",
			"last_name",
		}
		err := fmt.Errorf("%s %s", errorTexts[0], errorTexts[1])
		return errorTexts, err
	}

	return [3]string{}, nil
}

func createUser(ctx context.Context, user internalTools.User) (string, error) {

	// make DynamoDB map
	dynamoUserMap, err := attributevalue.MarshalMap(user)
	if err != nil {

		return "Error al convertir los datos a formato de DynamoDB.", err
	}

	// make input
	input := &dynamodb.PutItemInput{
		Item:      dynamoUserMap,
		TableName: aws.String(usersTableName),
	}

	// save
	_, err = dbSvc.PutItem(ctx, input)
	if err != nil {

		return "Error al intentar guardar el usuario.", err
	}

	return "", nil
}

func generateJWT(email string) (string, error) {
	claims := jwt.MapClaims{

		"email": email,

		"exp": time.Now().Add(time.Minute * 15).Unix(), // expiration in 15 min
		"iat": time.Now().Unix(),                       // emition time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func sendConfirmationToEmail(ctx context.Context, user internalTools.User, url string) error {

	input := &ses.SendEmailInput{
		Destination: &sesTypes.Destination{
			ToAddresses: []string{user.Email},
		},
		Message: &sesTypes.Message{
			Body: &sesTypes.Body{
				Html: &sesTypes.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(fmt.Sprintf("<html><body><h1>¡Hola desde SES!</h1><p>Este es un correo de prueba enviado con Go.</p><a>%s</a><p>(valido por 15 minutos)</p></body></html>", url)),
				},
				Text: &sesTypes.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(fmt.Sprintf("¡Hola desde SES!\nEste es un correo de prueba enviado con Go.\n%s\n(valido por 15 minutos)", url)),
				},
			},
			Subject: &sesTypes.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Correo de prueba con SES y Go"),
			},
		},
		Source: aws.String(senderEmailAddress), // Reemplaza con tu dirección de correo verificada en SES
	}

	if _, err := sesSvc.SendEmail(ctx, input); err != nil {
		return err
	}

	return nil
}

func updatePassword(ctx context.Context, email string, password string) (string, error) {

	// key
	key, err := attributevalue.MarshalMap(map[string]interface{}{
		"email": email,
	})
	if err != nil {

		return "Error al convertir los datos a formato de DynamoDB.", err
	}

	// make input
	input := &dynamodb.UpdateItemInput{
		TableName:        aws.String(usersTableName),
		Key:              key,
		UpdateExpression: aws.String("set #p = :val"),
		ExpressionAttributeNames: map[string]string{
			"#p": "password",
		},
		ExpressionAttributeValues: map[string]dynamodbTypes.AttributeValue{
			":val": &dynamodbTypes.AttributeValueMemberS{Value: password},
		},
	}

	// update
	_, err = dbSvc.UpdateItem(ctx, input)
	if err != nil {
		return "Error al actualizar la contraseña del nuevo usuario.", err
	}

	return "", nil
}

func getUser(ctx context.Context, email string) (*smallUser, error) {

	key, err := attributevalue.MarshalMap(map[string]interface{}{
		"email": email,
	})

	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName:            aws.String(usersTableName),
		Key:                  key,
		ConsistentRead:       aws.Bool(true),
		ProjectionExpression: aws.String("email,password,verified_email"),
	}

	result, err := dbSvc.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var user smallUser
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
