package models

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/pusher/pusher-http-go"
)

// EventPublisher - the object that can publish events
type EventPublisher struct {
	snsSession *sns.SNS
	topicARN   *string
	sqsSession *sqs.SQS
	queueARN   *string
}

// NewTopicEventPublisher - Create a new event publisher
func NewTopicEventPublisher(arn string) *EventPublisher {
	return &EventPublisher{
		snsSession: sns.New(session.Must(session.NewSession())),
		topicARN:   aws.String(arn),
	}
}

// PublishTopicEvent - publish event
func (eventPublisher *EventPublisher) PublishTopicEvent(conversation *Conversation, debug bool) error {
	params := &sns.PublishInput{
		Message:  aws.String(conversation.ToJSON()),
		TopicArn: eventPublisher.topicARN,
	}

	_, err := eventPublisher.snsSession.Publish(params)
	return err
}

// PublishMessageTopicEvent - publish event containing message information
func (eventPublisher *EventPublisher) PublishMessageTopicEvent(body string, guestPhone string, companyPhone string, conversationID string, debug bool) error {
	message, parseErr := MessageEventToJSON(body, guestPhone, companyPhone, conversationID)
	if parseErr != nil {
		log.Print("Unable to create message event json")
		log.Print(parseErr)
	}

	if debug {
		log.Printf("[PublishMessageTopicEvent] message: %s", message)
	}

	params := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: eventPublisher.topicARN,
	}

	_, err := eventPublisher.snsSession.Publish(params)
	return err
}

// PublishMessageQueueEvent - publish event containing message information to given queue
func (eventPublisher *EventPublisher) PublishMessageQueueEvent(body string, guestPhone string, companyPhone string, conversationID string, debug bool) error {
	message, parseErr := MessageEventToJSON(body, guestPhone, companyPhone, conversationID)
	if parseErr != nil {
		log.Print("Unable to create message event json")
		log.Print(parseErr)
	}

	params := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    eventPublisher.queueARN,
	}

	if debug {
		log.Printf("[PublishMessageQueueEvent] message: %s", message)
	}

	_, err := eventPublisher.sqsSession.SendMessage(params)
	return err
}

// PublishTopicMessage - publish topic message
func (eventPublisher *EventPublisher) PublishTopicMessage(body string) error {
	params := &sns.PublishInput{
		Message:  aws.String(body),
		TopicArn: eventPublisher.topicARN,
	}

	_, err := eventPublisher.snsSession.Publish(params)
	return err
}

// NewQueueEventPublisher - Create a new event publisher
func NewQueueEventPublisher(queueARN string) *EventPublisher {
	return &EventPublisher{
		sqsSession: sqs.New(session.Must(session.NewSession())),
		queueARN:   aws.String(queueARN),
	}
}

// PublishMessageFeedEvent - Publish a new message feed
func PublishMessageFeedEvent(appID string, key string, secret string, cluster string, data map[string]string, channelNames []string, eventName string) error {
	pusherClient := pusher.Client{
		AppID:   appID,
		Key:     key,
		Secret:  secret,
		Cluster: cluster,
		Secure:  true,
	}

	return pusherClient.TriggerMulti(channelNames, eventName, data)
}

// AuthenticatePrivateChannel - Authenticate websocket connection request
func AuthenticatePrivateChannel(appID string, key string, secret string, cluster string, params []byte) ([]byte, error) {
	pusherClient := pusher.Client{
		AppID:   appID,
		Key:     key,
		Secret:  secret,
		Cluster: cluster,
		Secure:  true,
	}
	return pusherClient.AuthenticatePrivateChannel(params)
}

// PublishConversationQueueEvent - publish event
func (eventPublisher *EventPublisher) PublishConversationQueueEvent(conversation *Conversation, debug bool) error {
	params := &sqs.SendMessageInput{
		MessageBody: aws.String(conversation.ToJSON()),
		QueueUrl:    eventPublisher.queueARN,
	}

	_, err := eventPublisher.sqsSession.SendMessage(params)
	return err
}

// PublishConversationMessageQueueEvent - publish conversation message event
func (eventPublisher *EventPublisher) PublishConversationMessageQueueEvent(conversationMessage *ConversationMessage, debug bool) error {
	message, parseErr := conversationMessage.ConvertToJSONObject().ToJSON()
	if parseErr != nil {
		log.Print("Error while publishing conversation message")
		log.Print(parseErr)
		return parseErr
	}

	params := &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    eventPublisher.queueARN,
	}

	_, err := eventPublisher.sqsSession.SendMessage(params)
	return err
}
