package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/nimanozari/restful-on-aws/lib"
	"testing"
)

type mockDB struct {
	dynamodbiface.DynamoDBAPI
}

func (m *mockDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, nil
}

type testPair struct {
	input          deviceLib.Request
	expectedOutput deviceLib.Response
}

func prepareTestPairs(pairs *[]testPair) {
	var myTestPair testPair

	myTestPair = testPair{
		input: deviceLib.Request{
			Body: *aws.String(`{"id":"/devices/id1","deviceModel":"/devicemodels/id1","name":"Sensor","note":"Testing a sensor.","serial":"A020000102"}`),
		},
		expectedOutput: deviceLib.Response_StatusCreated(`{"id":"/devices/id1","deviceModel":"/devicemodels/id1","name":"Sensor","note":"Testing a sensor.","serial":"A020000102"}`),
	}
	*pairs = append(*pairs, myTestPair)

	myTestPair = testPair{
		input: deviceLib.Request{
			Body: *aws.String(`{"id":"/devices/id1","model":"/devicemodels/id1","name":"Sensor","note":"Testing a sensor.","serial":"A020000102"}`),
		},
		expectedOutput: deviceLib.Response_BadRequest(`,"Message":"The following field is missing: deviceModel"`),
	}
	*pairs = append(*pairs, myTestPair)

	myTestPair = testPair{
		input: deviceLib.Request{
			Body: *aws.String(`{}`),
		},
		expectedOutput: deviceLib.Response_BadRequest(`,"Message":"The following fields are missing: id,deviceModel,name,note,serial"`),
	}
	*pairs = append(*pairs, myTestPair)

	myTestPair = testPair{
		input: deviceLib.Request{
			Body: *aws.String("Hello!"),
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

func TestInsertIntoDB(t *testing.T) {
	db := &mockDB{}
	var tests []testPair
	prepareTestPairs(&tests)
	for _, pair := range tests {
		v, _ := InsertIntoDB(pair.input, db)
		if !equal(v, pair.expectedOutput) {
			t.Error(
				"For", pair.input,
				"expected", pair.expectedOutput,
				"got", v,
			)
		}
	}
}
