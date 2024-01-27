package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

// EventMessage - Login params
type EventMessage struct {
	FromPhoneNumber string `json:"fromPhoneNumber"`
	ToPhoneNumber   string `json:"toPhoneNumber"`
	Body            string `json:"body"`
	CreatedAt       string `json:"createdAt"`
	ConversationID  string `json:"conversation_id"`
}

// ParseEventMessage - parse the sns event
func ParseEventMessage(snsEvent events.SNSEvent) []EventMessage {
	var messages []EventMessage
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		message := &EventMessage{
			FromPhoneNumber: "",
			ToPhoneNumber:   "",
			Body:            "",
			CreatedAt:       "",
			ConversationID:  "",
		}
		parseErr := json.Unmarshal([]byte(snsRecord.Message), &message)
		if parseErr == nil {
			messages = append(messages, *message)
		} else {
			log.Print("FAILURE, need to put this into a DLQ")
			fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)
		}
	}

	return messages
}

// ToJSON - convert EventMessage to json string
func (message *EventMessage) ToJSON() (string, error) {
	mess, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(mess), nil
}
