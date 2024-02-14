package main

import (
	"log"
	"net/url"

	"github.com/Amo-Addai/devsecops-ci-cd/serverless-aws/models"
	"github.com/Amo-Addai/devsecops-ci-cd/serverless-aws/new_sms/lib"
	"github.com/Amo-Addai/devsecops-ci-cd/serverless-aws/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response - AWS Response
type Response events.APIGatewayProxyResponse

// Handler - main handler
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	environment := lib.NewEnvironment()
	if environment.Debug {
		log.Printf("topicURL: %s", environment.TopicURL)
	}

	data, err := url.ParseQuery(request.Body)
	if err != nil {
		log.Print("Unable to parse body")
		return Response{
			StatusCode: 503,
		}, nil
	}

	// validate params
	fromPhoneNumber := data.Get("From")
	toPhoneNumber := data.Get("To")
	body := data.Get("Body")
	if fromPhoneNumber == "" {
		log.Print("Missing From")
		return Response{
			StatusCode: 403,
		}, nil
	}

	// validate webhook
	twilio := lib.NewTexter(fromPhoneNumber, environment.TwilioAuthToken, environment.TwilioAccountSid, environment.ThisURL, data, environment.Debug)
	if !twilio.ValidateWebhook(request) {
		log.Print("Unable to validate webhook")
		return Response{
			StatusCode:      200,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/xml",
			},
		}, nil
	}

	// debug
	if environment.Debug {
		log.Printf("fromPhoneNumber: %s", fromPhoneNumber)
		log.Printf("toPhoneNumber: %s", toPhoneNumber)
		log.Printf("body: %s", body)
	}

	// send to topic
	message := &models.EventMessage{
		FromPhoneNumber: fromPhoneNumber,
		ToPhoneNumber:   toPhoneNumber,
		Body:            body,
		CreatedAt:       utils.Now(),
	}
	publisher := models.NewTopicEventPublisher(environment.TopicURL)
	event, jsonErr := message.ToJSON()
	if jsonErr != nil {
		log.Print("Unable to format message event to string")
	}
	publisher.PublishTopicMessage(event)

	// success
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
	}

	log.Print("Processed!")
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
