package main

import (
	// base

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func updateVerifiedEmail(ctx context.Context, email string) (string, error) {

	// key
	key, err := attributevalue.MarshalMap(map[string]interface{}{
		"email": email,
	})
	if err != nil {

		return "Error al construir la expresión para actualizar el usuario.", err
	}

	// make input
	input := &dynamodb.UpdateItemInput{
		TableName:        aws.String(usersTableName),
		Key:              key,
		UpdateExpression: aws.String("set #ve = :val"),
		ExpressionAttributeNames: map[string]string{
			"#ve": "verified_email",
		},
		ExpressionAttributeValues: map[string]dynamodbTypes.AttributeValue{
			":val": &dynamodbTypes.AttributeValueMemberBOOL{Value: true},
		},
	}

	// update
	_, err = dbSvc.UpdateItem(ctx, input)
	if err != nil {
		return "Error al actualizar el usuario.", err
	}

	return "", nil
}
