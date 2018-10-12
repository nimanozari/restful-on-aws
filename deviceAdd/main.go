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
	"strings"
)

func InsertIntoDB(req deviceLib.Request, db dynamodbiface.DynamoDBAPI) (deviceLib.Response, error) {
	var dev deviceLib.Device
	err := json.Unmarshal([]byte(req.Body), &dev)
	if err != nil {
		return deviceLib.Response_InternalServerError(""), nil
	}
	var missingFields string = ""
	var missingCount int = 0
	if dev.Id == "" {
		missingFields += `id,`
		missingCount++
	}
	if dev.DeviceModel == "" {
		missingFields += `deviceModel,`
		missingCount++
	}
	if dev.Name == "" {
		missingFields += `name,`
		missingCount++
	}
	if dev.Note == "" {
		missingFields += `note,`
		missingCount++
	}
	if dev.Serial == "" {
		missingFields += `serial,`
		missingCount++
	}
	if missingFields != "" {
		missingFields = strings.TrimSuffix(missingFields, ",")
		if missingCount == 1 {
			return deviceLib.Response_BadRequest(`,"Message":"The following field is missing: ` + missingFields + `"`), nil
		} else {
			return deviceLib.Response_BadRequest(`,"Message":"The following fields are missing: ` + missingFields + `"`), nil
		}
	}

	deviceData, err := dynamodbattribute.MarshalMap(dev)
	if err != nil {
		return deviceLib.Response_InternalServerError(""), nil
	}

	insertCommand := &dynamodb.PutItemInput{
		Item:      deviceData,
		TableName: aws.String("devicesDynamoDBTable"),
	}

	_, err = db.PutItem(insertCommand)
	if err != nil {
		return deviceLib.Response_InternalServerError(""), nil
	}

	js, err := json.Marshal(dev)
	if err != nil {
		return deviceLib.Response_InternalServerError(""), nil
	}

	return deviceLib.Response_StatusCreated(string(js)), nil
}

func Handler(req deviceLib.Request) (deviceLib.Response, error) {
	region := os.Getenv("AWS_REGION")
	mySession, sessionError := session.NewSession(&aws.Config{Region: &region})
	if sessionError != nil {
		return deviceLib.Response_InternalServerError(""), nil
	}
	db := dynamodb.New(mySession)

	return InsertIntoDB(req, db)
}

func main() {
	lambda.Start(Handler)
}
