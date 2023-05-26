package examples

import (
	"context"
	"testing"

	dynamock "github.com/alext234/go-dynamock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var mock *dynamock.DynaMock

func init() {
	Dyna = new(MyDynamo)
	Dyna.Db, mock = dynamock.New()
}

func TestGetName(t *testing.T) {
	expectKey := map[string]*dynamodb.AttributeValue{
		"id": {
			N: aws.String("1"),
		},
	}

	expectedResult := aws.String("jaka")
	result := dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"name": {
				S: expectedResult,
			},
		},
	}

	//lets start dynamock in action
	mock.ExpectGetItem().ToTable("employee").WithKeys(expectKey).WillReturns(result)

	actualResult, _ := GetName("1")
	if actualResult != expectedResult {
		t.Errorf("Test Fail")
	}
}

func TestGetTransactGetItems(t *testing.T) {
	databaseOutput := dynamodb.TransactWriteItemsOutput{}

	mock.ExpectTransactWriteItems().Table("wrongTable").WillReturns(databaseOutput)

	err := GetTransactGetItems("")

	if err == nil {
		t.Errorf("Test failed")
	}
}

func TestQueryItems(t *testing.T) {
	expectedResult := aws.String("haha")
	result := dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"name": {
					S: expectedResult,
				},
			},
		},
	}
	idValue := "1"
	expectedQuery := &dynamodb.QueryInput{
		TableName:                 aws.String("employee"),
		KeyConditionExpression:    aws.String("id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":id": {S: aws.String(idValue)}},
	}

	mock.ExpectQuery().Table("employee").WithQueryInput(expectedQuery).WillReturns(result)

	actualResult, err := QueryItems(idValue)
	if err != nil {
		t.Errorf("err is not nil %+v", err)
	}

	if actualResult != expectedResult {
		t.Errorf("Test Fail")
	}
}

func TestQueryItemsWithContext(t *testing.T) {
	expectedResult := aws.String("haha")
	result := dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"name": {
					S: expectedResult,
				},
			},
		},
	}
	idValue := "1"
	expectedQuery := &dynamodb.QueryInput{
		TableName:                 aws.String("employee"),
		KeyConditionExpression:    aws.String("id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":id": {S: aws.String(idValue)}},
	}

	ctx := context.Background()

	mock.ExpectQuery().Table("employee").WithContext(ctx).
		WithQueryInput(expectedQuery).
		WillReturns(result)

	actualResult, err := QueryItemsWithContext(ctx, idValue)
	if err != nil {
		t.Errorf("err is not nil %+v", err)
	}

	if actualResult != expectedResult {
		t.Errorf("Test Fail")
	}
}
