package utils

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// CreateSuccessfulResponse - Create successful response
func CreateSuccessfulResponse() (events.APIGatewayProxyResponse, error) {
	return createResponse("successful", 200)
}

// CreateSuccessfulResponseWithPayload - Create successful response
func CreateSuccessfulResponseWithPayload(payload map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	return createResponseWithPayload("successful", 200, payload)
}

// CreateUnauthorizedResponse - Create unauthorized response
func CreateUnauthorizedResponse() (events.APIGatewayProxyResponse, error) {
	return createResponse("unauthorized", 401)
}

// CreateFailureResponse - Create failure response
func CreateFailureResponse() (events.APIGatewayProxyResponse, error) {
	return createResponse("unavailable", 503)
}

// CreateServerErrorResponse - Create server error response
func CreateServerErrorResponse() (events.APIGatewayProxyResponse, error) {
	return createResponse("server error", 500)
}

// CreateBadRequestResponse - Create bad request response
func CreateBadRequestResponse() (events.APIGatewayProxyResponse, error) {
	return createResponse("bad request", 400)
}

func createResponse(message string, statusCode int) (events.APIGatewayProxyResponse, error) {
	var buf bytes.Buffer
	body, err := json.Marshal(map[string]interface{}{"message": message})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}
	json.HTMLEscape(&buf, body)
	resp := events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}

	return resp, nil
}

func createResponseWithPayload(message string, statusCode int, payload map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	var buf bytes.Buffer
	body, err := json.Marshal(payload)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}
	json.HTMLEscape(&buf, body)
	resp := events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}

	return resp, nil
}
