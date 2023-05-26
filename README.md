# go-dynamock
Amazon Dynamo DB Mock Driver for Golang to Test Database Interactions.

Originally forked from https://github.com/gusaul/go-dynamock and added more functionalities with examples:

* `QueryInput` expectation.
* `WithContext` expectation to support `QueryWithContext`.

## Install
```
go get github.com/alext234/go-dynamock
```

## Examples Usage
See the `examples` directory.

### DynamoDB configuration
First of all, change the dynamodb configuration to use the ***dynamodb interface***. see code below:
``` go
package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MyDynamo struct {
    Db dynamodbiface.DynamoDBAPI
}

var Dyna *MyDynamo

func ConfigureDynamoDB() {
	Dyna = new(MyDynamo)
	awsSession, _ := session.NewSession(&aws.Config{Region: aws.String("ap-southeast-2")})
	svc := dynamodb.New(awsSession)
	Dyna.Db = dynamodbiface.DynamoDBAPI(svc)
}
```
the purpose of code above is to make your dynamoDB object can be mocked by ***dynamock*** through the dynamodbiface.

### Something you may wanna test
``` go
package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
)

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
```

### Test with DynaMock
``` go
package examples

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	dynamock "github.com/alext234/go-dynamock"
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
```
if you just wanna expect the table
``` go
mock.ExpectGetItem().ToTable("employee").WillReturns(result)
```
or maybe you didn't care with any arguments, you just need to determine the result
``` go
mock.ExpectGetItem().WillReturns(result)
```
and you can do multiple expectations at once, then the expectation will be executed sequentially.
``` go
mock.ExpectGetItem().WillReturns(resultOne)
mock.ExpectUpdateItem().WillReturns(resultTwo)
mock.ExpectGetItem().WillReturns(resultThree)

/* Result
the first call of GetItem will return resultOne
the second call of GetItem will return resultThree
and the only call of UpdateItem will return resultTwo */
```

## License

The [MIT License](https://github.com/alext234/go-dynamock/blob/master/LICENSE)
