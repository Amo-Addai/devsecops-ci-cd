package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
)

// ConversationMessage - the conversation message
type ConversationMessage struct {
	ConversationID string `dynamo:"conversation_id"`
	CreatedAt      string `dynamo:"created_at"`
	SenderGUID     string `dynamo:"sender_guid"`
	RecipientGUID  string `dynamo:"recipient_guid"`
	Body           string `dynamo:"body"`
}

// ConversationMessageJSON - the conversation message representable in json
type ConversationMessageJSON struct {
	ConversationID string `json:"conversation_id"`
	CreatedAt      string `json:"created_at"`
	SenderGUID     string `json:"sender_guid"`
	RecipientGUID  string `json:"recipient_guid"`
	Body           string `json:"body"`
}

// MessageEvent - message when sent as an event
type messageEvent struct {
	Body           string `json:"body"`
	GuestPhone     string `json:"guest_phone"`
	CompanyPhone   string `json:"company_phone"`
	ConversationID string `json:"conversation_id"`
}

// ConvertToJSONObject - convert the conversationMessage to a json representable object
func (message *ConversationMessage) ConvertToJSONObject() *ConversationMessageJSON {
	return &ConversationMessageJSON{
		ConversationID: message.ConversationID,
		CreatedAt:      message.CreatedAt,
		SenderGUID:     message.SenderGUID,
		RecipientGUID:  message.RecipientGUID,
		Body:           message.Body,
	}
}

// AddConversationMessage - add conversation message, sending to guest
func AddConversationMessage(table dynamo.Table, now string, conversation *Conversation, body string) (*ConversationMessage, error) {
	conversationMessage := ConversationMessage{
		ConversationID: conversation.ConversationID,
		CreatedAt:      now,
		SenderGUID:     conversation.DestinationGUID,
		RecipientGUID:  conversation.GuestGUID,
		Body:           body,
	}
	return &conversationMessage, table.Put(conversationMessage).Run()
}

// AddConversationMessageFromGuest - add conversation message, from guest
func AddConversationMessageFromGuest(table dynamo.Table, now string, conversation *Conversation, body string) (*ConversationMessage, error) {
	conversationMessage := ConversationMessage{
		ConversationID: conversation.ConversationID,
		CreatedAt:      now,
		SenderGUID:     conversation.GuestGUID,
		RecipientGUID:  conversation.DestinationGUID,
		Body:           body,
	}
	return &conversationMessage, table.Put(conversationMessage).Run()
}

// MessageEventToJSON - get the message event as a json
func MessageEventToJSON(body string, guestPhone string, companyPhone string, conversationID string) (string, error) {
	obj := &messageEvent{
		Body:           body,
		GuestPhone:     guestPhone,
		CompanyPhone:   companyPhone,
		ConversationID: conversationID,
	}
	convo, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(convo), nil
}

// ToJSON - return ConversationMessageJSON as json string
func (message *ConversationMessageJSON) ToJSON() (string, error) {
	mess, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(mess), nil
}

// FromJSONToConversationMessageJSON - convert string to ConversationMessageJSON
func FromJSONToConversationMessageJSON(jsonString string) (*ConversationMessageJSON, error) {
	message := &ConversationMessageJSON{
		Body:           "",
		ConversationID: "",
		CreatedAt:      "",
		SenderGUID:     "",
		RecipientGUID:  "",
	}

	parseErr := json.Unmarshal([]byte(jsonString), &message)
	if parseErr != nil {
		log.Print("FAILURE, unable to unmarshall string to ConversationMessageJSON")
		log.Printf("Message = %s \n", jsonString)
		return nil, parseErr
	}
	return message, nil
}

// MessageEventFromJSON - get the message event as a JSON
func MessageEventFromJSON(jsonString string) (string, string, string, string, error) {
	message := &messageEvent{
		Body:           "",
		GuestPhone:     "",
		CompanyPhone:   "",
		ConversationID: "",
	}

	parseErr := json.Unmarshal([]byte(jsonString), &message)
	if parseErr != nil {
		log.Print("FAILURE, need to put this into a DLQ")
		log.Printf("Message = %s \n", jsonString)
	}
	return message.Body, message.GuestPhone, message.CompanyPhone, message.ConversationID, parseErr
}

// ParseMessageEventFromQueue - parse the sqs event to a event message
//   A bit confusing because MessageEvent is an internal representation of a message.
//   But an EventMessage is an object represents an external message event.
func ParseMessageEventFromQueue(sqsEvent events.SQSEvent) []EventMessage {
	var messages []EventMessage
	for _, record := range sqsEvent.Records {

		body, guestPhone, companyPhone, conversationID, err := MessageEventFromJSON(record.Body)
		if err != nil {
			log.Print("FAILURE, need to put this into a DLQ")
			log.Printf("[ParseMessageEventFromQueue] - Message = %s \n", record.Body)
			continue
		}

		message := &EventMessage{
			FromPhoneNumber: guestPhone,
			ToPhoneNumber:   companyPhone,
			Body:            body,
			ConversationID:  conversationID,
		}
		messages = append(messages, *message)
	}

	return messages
}

// GetConversationMessagesByConversationGUID - Gets all the conversation messages by the conversation guid.
func GetConversationMessagesByConversationGUID(table dynamo.Table, conversationGUID string, debug bool) ([]ConversationMessage, error) {
	var errToReturn error
	var items []ConversationMessage

	if debug {
		fmt.Printf("[GetConversationMessagesByConversationGUID] - Start conversation guid: %s", conversationGUID)
	}
	err := table.Get("conversation_id", conversationGUID).All(&items)
	if err != nil {
		if err == dynamo.ErrNotFound {
			if debug {
				fmt.Printf("[GetConversationMessagesByConversationGUID] No messages found for conversation guid %s, from orm", conversationGUID)
			}
			errToReturn = nil
		} else {
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
					if debug {
						fmt.Printf("[GetConversationMessagesByConversationGUID] No messages found for conversation guid %s, from dynamodb", conversationGUID)
					}
					errToReturn = nil
				} else {
					if debug {
						fmt.Printf("[GetConversationMessagesByConversationGUID] Bad aerr for conversation guid %s", conversationGUID)
						log.Print(aerr)
					}
					errToReturn = aerr
				}
			} else {
				fmt.Printf("[GetConversationMessagesByConversationGUID] Error while trying to find messages for conversationGUID: %s", conversationGUID)
				log.Print(err)
				errToReturn = err
			}
		}
		return items, errToReturn
	}
	if debug {
		fmt.Printf("[GetConversationMessagesByConversationGUID] End - Messages found for conversation guid: %s!", conversationGUID)
	}
	return items, nil
}
