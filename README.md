# Simple Restful API on AWS

A simple Restful API on AWS implemented using the following tech stack:

* [Serverless Framework](https://serverless.com/)
* [Go language](https://golang.org/)
* AWS API Gateway
* AWS Lambda
* AWS DynamoDB



## Specifications:

The API accepts the following JSON requests and produces the corresponding HTTP responses:


### Request 1:
`HTTP POST
URL: https://<api-gateway-url>/api/devices`

#### Response 1 - Success:
HTTP 201 Created
Body (application/json):
{
  "id": "/devices/id1",
  "deviceModel": "/devicemodels/id1",
  "name": "Sensor",
  "note": "Testing a sensor.",
  "serial": "A020000102"
}

#### Response 1 - Failure 1:
HTTP 400 Bad Request
If any of the payload fields are missing, response body will include a descriptive error message for the client to be able to detect the problem.

#### Response 1 - Failure 2:
HTTP 500 Internal Server Error
If any exceptional situation occurs on the server side.


### Request 2:
HTTP GET

URL: https://<api-gateway-url>/api/devices/{id}

Example: GET https://api123.amazonaws.com/api/devices/id1

#### Response 2 - Success:
HTTP 200 OK
Body (application/json):
{
  "id": "/devices/id1",
  "deviceModel": "/devicemodels/id1",
  "name": "Sensor",
  "note": "Testing a sensor.",
  "serial": "A020000102"
}

#### Response 2 - Failure 1:
HTTP 404 Not Found
If the request id does not exist.

#### Response 2 - Failure 2:
HTTP 500 Internal Server Error
If any exceptional situation occurs on the server side.


## How To Run & Test:

After successfully setting up the Go Language and Serverless Framework on the system, and setting the Amazon Web Services credentials, the code can be deployed by navigating to the directory containing the project and typing the following command in the terminal:

`make deploy`

If done successfully, it will print some URLs to access the corresponding servicess inside the terminal.

To test the POST method, the `curl` command can be used from the terminal:

`curl -i -H "Content-Type: application/json" -X POST [THE CORRESPONDING URL ADDRESS] -d '{"id":"/devices/id1","deviceModel":"/devicemodels/id1","name":"Sensor","note":"Testing a sensor.","serial":"A020000102"}'`

To test the GET method, simply visit the given URL using any browser, and replace the trailing `{id}` with `id1`.

### Unit Test Coverage:

A sample unit test is included for both POST and GET functions in their respective directories, which can be utilized by navigating to the directory and entering the following commands in the terminal:

`go test -coverprofile=coverage.out`

To see the test unit coverage percentages:

`go tool cover -func=coverage.out`

To see which parts of code are covered:

`go tool cover -html=coverage.out`
