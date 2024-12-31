package main

import (
	// base

	"context"
	"time"

	// from third parties

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/golang-jwt/jwt/v5"

	"github.com/GuidoGdR/users-GO-AWS/usersTools"
)

func generateJWT(user *smallUser) (string, error) {
	claims := usersTools.JwtClaims{
		Email: user.Email,

		FirstName: user.FirstName,
		LastName:  user.LastName,

		Exp: time.Now().Add(time.Hour * 24).Unix(),
		Iat: time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
		ProjectionExpression: aws.String("email,password,first_name,last_name,email_verified"),
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
