package deviceLib

import (
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

type Device struct {
	Id          string `json:"id"`
	DeviceModel string `json:"deviceModel"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Serial      string `json:"serial"`
}

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

func Response_InternalServerError(body string) Response {
	return Response{
		StatusCode: http.StatusInternalServerError,
		Body:       `{"Error":"HTTP 500 Internal Server Error."` + body + "}",
	}
}

func Response_BadRequest(body string) Response {
	return Response{
		StatusCode: http.StatusBadRequest,
		Body:       `{"Error":"HTTP 400 Bad Request."` + body + "}",
	}
}

func Response_NotFound(body string) Response {
	return Response{
		StatusCode: http.StatusNotFound,
		Body:       `{"Error":"HTTP 404 Not Found."` + body + "}",
	}
}

func Response_StatusCreated(body string) Response {
	return Response{
		StatusCode: http.StatusCreated,
		Body:       body,
	}
}

func Response_OK(body string) Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       body,
	}
}
