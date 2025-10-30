package main

import (
	// base

	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func updateVerifiedEmail(ctx context.Context, email string) (bool, string, error) {

	// key
	key, err := attributevalue.MarshalMap(map[string]interface{}{
		"email": email,
	})
	if err != nil {

		return false, "Error al construir la expresión para actualizar el usuario.", err
	}

	// make input
	input := &dynamodb.UpdateItemInput{
		TableName:        aws.String(usersTableName),
		Key:              key,
		UpdateExpression: aws.String("set #ve = :t"),
		ExpressionAttributeNames: map[string]string{
			"#ve": "verified_email",
		},
		ExpressionAttributeValues: map[string]dynamodbTypes.AttributeValue{
			":t": &dynamodbTypes.AttributeValueMemberBOOL{Value: true},
			":f": &dynamodbTypes.AttributeValueMemberBOOL{Value: false},
		},
		ConditionExpression: aws.String("#ve = :f"),
	}

	// update
	_, err = dbSvc.UpdateItem(ctx, input)
	if err != nil {

		//2
		var conditionalCheckFailedException *dynamodbTypes.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailedException) {

			return true, "El usuario no existe o el email ya ha sido verificado.", err // user dont exist
		}

		return false, "Error al actualizar el usuario.", err
	}

	return false, "", nil
}
