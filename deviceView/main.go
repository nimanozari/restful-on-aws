package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/nimanozari/restful-on-aws/lib"
	"os"
)

func ReadFromDB(req deviceLib.Request, db dynamodbiface.DynamoDBAPI) (deviceLib.Response, error) {
	key := map[string]*dynamodb.AttributeValue{
		"id": {S: aws.String(req.Path)},
	}
	readCommand := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String("devicesDynamoDBTable"),
	}
	var dev deviceLib.Device
	result, err := db.GetItem(readCommand)
	if err != nil {
		return deviceLib.Response_InternalServerError(""), nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &dev)
	if err != nil {
		return deviceLib.Response_InternalServerError(""), nil
	}
	if dev.Id == "" {
		return deviceLib.Response_NotFound(`,"Message":"Unable to find an item with such id."`), nil
	}

	js, err := json.Marshal(dev)
	if err != nil {
		return deviceLib.Response_InternalServerError(""), nil
	}

	return deviceLib.Response_OK(string(js)), nil
}

func Handler(req deviceLib.Request) (deviceLib.Response, error) {
	region := os.Getenv("AWS_REGION")
	mySession, sessionError := session.NewSession(&aws.Config{Region: &region})
	if sessionError != nil {
		return deviceLib.Response_InternalServerError(""), nil
	}
	db := dynamodb.New(mySession)

	return ReadFromDB(req, db)
}

func main() {
	lambda.Start(Handler)
}
