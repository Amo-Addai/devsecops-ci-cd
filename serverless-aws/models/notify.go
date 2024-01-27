package models

import (
	"log"

	"github.com/guregu/dynamo"
)

// Notify - Phone numbers to notify
type Notify struct {
	CompanyID  string `dynamo:"company_id"`
	PhoneAdmin string `dynamo:"phone_admin"`
	UserID     string `dynamo:"user_id"` // optional cognito user ids
}

// GetNotifiesForConversation - Get conversations for a given phone
func GetNotifiesForConversation(table dynamo.Table, companyID string) ([]Notify, error) {
	var items []Notify

	err := table.Scan().Filter("'company_id' = ?", companyID).All(&items)
	if err != nil {
		log.Println("Error while getting notifies")
		log.Println(err)
		return items, err
	}

	return items, nil
}
