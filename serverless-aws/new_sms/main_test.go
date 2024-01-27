package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestHandler(t *testing.T) {
	request := events.APIGatewayProxyRequest{}
	response, err := Handler(request)

	assert.Nil(t, err)
	assert.NotNil(t, response)
}
