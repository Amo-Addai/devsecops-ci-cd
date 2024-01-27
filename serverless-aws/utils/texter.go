package utils

import (
	"log"

	"github.com/sfreiberg/gotwilio"
)

// Texter - the object that can text
type Texter struct {
	twilio *gotwilio.Twilio
}

// NewTexter - creates a new texter
func NewTexter(accountSID string, authToken string) *Texter {
	return &Texter{
		twilio: gotwilio.NewTwilioClient(accountSID, authToken),
	}
}

// Text - text a phone number
func (texter *Texter) Text(from string, to string, message string) error {
	smsResponse, exception, err := texter.twilio.SendSMS(from, to, message, "", "")
	if err != nil {
		log.Print("Unable to send text, err received")
		log.Print(exception)
		log.Print(err)
	}
	if exception != nil {
		log.Print("Unable to send text, exception received")
		log.Print(exception)
		log.Print(err)
	}
	log.Print(smsResponse)
	return err
}
