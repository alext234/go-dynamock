package examples

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// MyDynamo struct hold dynamodb connection
type MyDynamo struct {
	Db dynamodbiface.DynamoDBAPI
}

// Dyna - object from MyDynamo
var Dyna *MyDynamo

// ConfigureDynamoDB - init func for open connection to aws dynamodb
func ConfigureDynamoDB() {
	Dyna = new(MyDynamo)
	awsSession, _ := session.NewSession(&aws.Config{Region: aws.String("ap-southeast-2")})
	svc := dynamodb.New(awsSession)
	Dyna.Db = dynamodbiface.DynamoDBAPI(svc)
}

// GetName - example func using GetItem method
func GetName(id string) (*string, error) {
	parameter := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
		TableName: aws.String("employee"),
	}

	response, err := Dyna.Db.GetItem(parameter)
	if err != nil {
		return nil, err
	}

	name := response.Item["name"].S
	return name, nil
}

// GetName - example func using GetItem method
func GetTransactGetItems(id string) error {
	parameter := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName: aws.String("my_table"),
				},
			},
		},
	}

	_, err := Dyna.Db.TransactWriteItems(parameter)

	if err != nil {
		fmt.Print(err.Error())
		return err
	}

	return nil
}

// QueryItems - example func using Query method
func QueryItems(id string) (*string, error) {
	query := &dynamodb.QueryInput{
		TableName:                 aws.String("employee"),
		KeyConditionExpression:    aws.String("id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":id": {S: aws.String(id)}},
	}

	result, err := Dyna.Db.Query(query)

	if err != nil {
		return nil, err
	}

	for _, item := range result.Items {
		return item["name"].S, nil
	}
	return nil, errors.New("empty result")
}

// QueryItems - example func using QueryWithContext method
func QueryItemsWithContext(ctx context.Context, id string) (*string, error) {
	query := &dynamodb.QueryInput{
		TableName:                 aws.String("employee"),
		KeyConditionExpression:    aws.String("id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":id": {S: aws.String(id)}},
	}

	result, err := Dyna.Db.QueryWithContext(ctx, query)

	if err != nil {
		return nil, err
	}

	for _, item := range result.Items {
		return item["name"].S, nil
	}
	return nil, errors.New("empty result")
}
