package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"nima-test/lib"
	"testing"
)

type mockDB struct {
	dynamodbiface.DynamoDBAPI
}

func (m *mockDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {

	if *input.Key["id"].S == "Something that causes database error" {
		return &dynamodb.GetItemOutput{}, errors.New("db error")
	} else if *input.Key["id"].S == "/devices/id1" {
		resultItem := map[string]*dynamodb.AttributeValue{
			"id":          {S: aws.String("/devices/id1")},
			"deviceModel": {S: aws.String("/devicemodels/id1")},
			"name":        {S: aws.String("Sensor")},
			"note":        {S: aws.String("Testing a sensor.")},
			"serial":      {S: aws.String("A020000102")},
		}
		return &dynamodb.GetItemOutput{Item: resultItem}, nil
	} else {
		return &dynamodb.GetItemOutput{}, nil
	}
}

type testPair struct {
	input          deviceLib.Request
	expectedOutput deviceLib.Response
}

func prepareTestPairs(pairs *[]testPair) {
	var myTestPair testPair
	myTestPair = testPair{
		input: deviceLib.Request{
			Path: *aws.String("/devices/id1"),
		},
		expectedOutput: deviceLib.Response_OK(`{"id":"/devices/id1","deviceModel":"/devicemodels/id1","name":"Sensor","note":"Testing a sensor.","serial":"A020000102"}`),
	}
	*pairs = append(*pairs, myTestPair)

	myTestPair = testPair{
		input: deviceLib.Request{
			Path: *aws.String("/devices/id88"),
		},
		expectedOutput: deviceLib.Response_NotFound(`,"Message":"Unable to find an item with such id."`),
	}
	*pairs = append(*pairs, myTestPair)

	myTestPair = testPair{
		input: deviceLib.Request{
			Path: *aws.String("Something that causes database error"),
		},
		expectedOutput: deviceLib.Response_InternalServerError(""),
	}
	*pairs = append(*pairs, myTestPair)
}

func equal(r1 deviceLib.Response, r2 deviceLib.Response) bool {
	if r1.StatusCode != r2.StatusCode {
		return false
	}
	if r1.Body != r2.Body {
		return false
	}
	return true
}

func TestReadFromDB(t *testing.T) {
	db := &mockDB{}
	var tests []testPair
	prepareTestPairs(&tests)
	for _, pair := range tests {
		v, _ := ReadFromDB(pair.input, db)
		if !equal(v, pair.expectedOutput) {
			t.Error(
				"For", pair.input,
				"expected", pair.expectedOutput,
				"got", v,
			)
		}
	}
}
