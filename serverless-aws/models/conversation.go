package models

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Amo-Addai/devsecops-ci-cd/serverless-aws/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
)

// Conversation - conversation
type Conversation struct {
	PhoneGuest           string `dynamo:"phone_guest"`
	PhoneDestination     string `dynamo:"phone_dest"`
	ConversationID       string `dynamo:"conversation_id"`
	GuestGUID            string `dynamo:"guest_guid"`
	DestinationGUID      string `dynamo:"dest_guid"`
	CreatedAt            string `dynamo:"created_at"`
	BotEnabled           bool   `dynamo:"bot_enabled"`
	BoostConversationID  string `dynamo:"boost_conversation_id"`
	GuestName            string `dynamo:"guest_name"`
	GuestNotes           string `dynamo:"guest_notes"`
	CheckedAt            string `dynamo:"checked_at"`
	Read                 bool   `dynamo:"read"`
	LastMessage          string `dynamo:"last_message"`
	LastMessageCreatedAt string `dynamo:"last_message_created_at"`
}

// ConversationJSON - Conversation that can be serialized
type ConversationJSON struct {
	PhoneGuest           string `json:"phone_guest"`
	PhoneDestination     string `json:"phone_dest"`
	ConversationID       string `json:"conversation_id"`
	GuestGUID            string `json:"guest_guid"`
	DestinationGUID      string `json:"dest_guid"`
	CreatedAt            string `json:"created_at"`
	BotEnabled           bool   `json:"bot_enabled"`
	BoostConversationID  string `json:"boost_conversation_id"`
	GuestName            string `json:"guest_name"`
	GuestNotes           string `json:"guest_notes"`
	CheckedAt            string `json:"checked_at"`
	Read                 bool   `json:"read"`
	LastMessage          string `json:"last_message"`
	LastMessageCreatedAt string `json:"last_message_created_at"`
}

// ToJSONObject - Convert to conversations that are capable to serializing to json
func (conversation *Conversation) ToJSONObject() *ConversationJSON {
	return &ConversationJSON{
		PhoneGuest:           conversation.PhoneGuest,
		PhoneDestination:     conversation.PhoneDestination,
		ConversationID:       conversation.ConversationID,
		GuestGUID:            conversation.GuestGUID,
		DestinationGUID:      conversation.DestinationGUID,
		CreatedAt:            conversation.CreatedAt,
		BotEnabled:           conversation.BotEnabled,
		BoostConversationID:  conversation.BoostConversationID,
		GuestName:            conversation.GuestName,
		GuestNotes:           conversation.GuestNotes,
		CheckedAt:            conversation.CheckedAt,
		Read:                 conversation.Read,
		LastMessage:          conversation.LastMessage,
		LastMessageCreatedAt: conversation.LastMessageCreatedAt,
	}
}

// ToJSON - Convert conversation to json
func (conversation *Conversation) ToJSON() string {
	convoJSON := conversation.ToJSONObject()
	convo, err := json.Marshal(convoJSON)
	if err != nil {
		log.Print("Unable to convert conversation to JSON")
		log.Print(conversation)
		return ""
	}
	return string(convo)
}

// GetGroupName - Get the group name for this conversation.
func (conversation *Conversation) GetGroupName() string {
	phoneRe := regexp.MustCompile(`^\+(\d+)$`)
	phoneMatch := phoneRe.FindStringSubmatch(conversation.PhoneDestination)
	phone := conversation.PhoneDestination
	if len(phoneMatch) == 2 {
		phone = phoneMatch[1]
	}
	return fmt.Sprintf("group_%s", phone)
}

// NewConversationFromJSON - convert json string to conversation
func NewConversationFromJSON(jsonString string) *Conversation {
	var conversation *Conversation
	convo := &ConversationJSON{
		PhoneGuest:           "",
		PhoneDestination:     "",
		ConversationID:       "",
		GuestGUID:            "",
		DestinationGUID:      "",
		CreatedAt:            "",
		BotEnabled:           true,
		BoostConversationID:  "",
		GuestName:            "",
		GuestNotes:           "",
		CheckedAt:            "",
		Read:                 false,
		LastMessage:          "",
		LastMessageCreatedAt: "",
	}

	parseErr := json.Unmarshal([]byte(jsonString), &convo)
	if parseErr == nil {
		conversation = &Conversation{
			PhoneGuest:           convo.PhoneGuest,
			PhoneDestination:     convo.PhoneDestination,
			ConversationID:       convo.ConversationID,
			GuestGUID:            convo.GuestGUID,
			DestinationGUID:      convo.DestinationGUID,
			CreatedAt:            convo.CreatedAt,
			BotEnabled:           convo.BotEnabled,
			BoostConversationID:  convo.BoostConversationID,
			GuestName:            convo.GuestName,
			GuestNotes:           convo.GuestNotes,
			CheckedAt:            convo.CheckedAt,
			Read:                 convo.Read,
			LastMessage:          convo.LastMessage,
			LastMessageCreatedAt: convo.LastMessageCreatedAt,
		}
	} else {
		log.Print("FAILURE, need to put this into a DLQ")
		log.Printf("Message = %s \n", jsonString)
	}
	return conversation
}

// FindConversation - find conversation
func FindConversation(table dynamo.Table, phoneGuest string, phoneDest string, debug bool) (*Conversation, error) {
	var errToReturn error
	var conversation *Conversation

	if debug {
		log.Printf("[FindConversation] phone_guest: %s, phone_dest: %s", phoneGuest, phoneDest)
	}
	err := table.Get("phone_guest", phoneGuest).
		Range("phone_dest", dynamo.Equal, phoneDest).
		One(&conversation)
	if err != nil {
		if err == dynamo.ErrNotFound {
			if debug {
				log.Print("[FindConversation] No convo found, from orm")
			}
			errToReturn = nil
		} else {
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
					if debug {
						log.Print("[FindConversation] No convo found, from dynamodb")
					}
					errToReturn = nil
				} else {
					if debug {
						log.Print("[FindConversation] Bad aerr")
						log.Print(aerr)
					}
					errToReturn = aerr
				}
			} else {
				log.Print("[FindConversation] Error while trying to find existing conversation")
				log.Print(err)
				errToReturn = err
			}
		}
		return conversation, errToReturn
	}
	if debug {
		log.Print("[FindConversation] Convo found!")
	}
	return conversation, nil
}

// GetConversationForPhones - Get conversations for given company phones (phone_dest)
func GetConversationForPhones(region string, tableName string, companyPhones []string, debug bool) ([]Conversation, error) {
	var items []Conversation
	if companyPhones == nil || len(companyPhones) == 0 {
		return items, nil
	}

	// setup
	awscfg := &aws.Config{}
	awscfg.WithRegion(region)
	sess := session.Must(session.NewSession(awscfg))
	svc := dynamodb.New(sess)

	// build expression
	filt := expression.Name("phone_dest").Equal(expression.Value(companyPhones[0]))
	companyPhonesLength := len(companyPhones)
	if companyPhonesLength > 1 {
		for i := 1; i < companyPhonesLength; i++ {
			phone := companyPhones[i]
			filt = filt.Or(expression.Name("phone_dest").Equal(expression.Value(phone)))
		}
	}

	// build
	expr, buildExpressionErr := expression.NewBuilder().WithFilter(filt).Build()
	if buildExpressionErr != nil {
		return items, buildExpressionErr
	}

	// create input
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 &tableName,
	}

	// scan
	result, err := svc.Scan(params)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return items, nil
		}
		log.Println("Error while getting conversation for company phones")
		return items, err
	}

	// parse result
	items = make([]Conversation, *result.Count)
	for i, item := range result.Items {
		items[i] = Conversation{
			PhoneGuest:           utils.ExtractDynamoResultString(item, "phone_guest"),
			PhoneDestination:     utils.ExtractDynamoResultString(item, "phone_dest"),
			ConversationID:       utils.ExtractDynamoResultString(item, "conversation_id"),
			GuestGUID:            utils.ExtractDynamoResultString(item, "guest_guid"),
			DestinationGUID:      utils.ExtractDynamoResultString(item, "dest_guid"),
			CreatedAt:            utils.ExtractDynamoResultString(item, "created_at"),
			BotEnabled:           utils.ExtractDynamoResultBool(item, "bot_enabled"),
			BoostConversationID:  utils.ExtractDynamoResultString(item, "boost_conversation_id"),
			GuestName:            utils.ExtractDynamoResultString(item, "guest_name"),
			GuestNotes:           utils.ExtractDynamoResultString(item, "guest_notes"),
			CheckedAt:            utils.ExtractDynamoResultString(item, "checked_at"),
			Read:                 utils.ExtractDynamoResultBool(item, "read"),
			LastMessage:          utils.ExtractDynamoResultString(item, "last_message"),
			LastMessageCreatedAt: utils.ExtractDynamoResultString(item, "last_message_created_at"),
		}
	}
	return items, nil
}

// GetConversationsSince - Gets all the conversations since time, specifically for last message creation time
func GetConversationsSince(table dynamo.Table, since time.Time, debug bool) ([]Conversation, error) {
	var items []Conversation

	err := table.Scan().Filter("last_message_created_at >= ?", utils.GetTimeAsString(since)).Index("conversation_timestamp_index").All(&items)
	if err != nil {
		log.Println("Error while getting conversations")
		log.Println(err)
		return items, err
	}

	return items, nil
}

// GetConversations - Gets all the conversations
// if limit is 0, means no limit
func GetConversations(table dynamo.Table, limit int64, debug bool) ([]Conversation, error) {
	var items []Conversation

	var err error = nil
	if limit <= 0 {
		err = table.Scan().All(&items)
	} else {
		_, err = table.Scan().SearchLimit(limit).AllWithLastEvaluatedKey(&items)
		// first param is the key used for future searches
		// to get next, where key is the first param we are currently ignoring
		// nextKey, err := table.Scan().StartFrom(key).SearchLimit(limit).AllWithLastEvaluatedKey(&items)
	}
	if err != nil {
		log.Println("Error while getting conversations")
		log.Println(err)
		return items, err
	}

	return items, nil
}

// GetConversationByID - Get conversation by id
func GetConversationByID(table dynamo.Table, conversationID string, debug bool) (*Conversation, error) {
	conversation := &Conversation{}
	err := table.Get("conversation_id", conversationID).Index("conversation_id_index").One(&conversation)
	if err != nil {
		log.Println("Error while getting conversation by id")
		log.Println(err)
		return conversation, err
	}
	return conversation, nil
}

// IsUserOnConversation - Indicates if the user is on the conversation
func IsUserOnConversation(conversation *Conversation, user *User, debug bool) bool {
	if conversation == nil {
		return false
	}

	if user.Superadmin {
		return true
	}

	for _, phone := range user.DestinationPhones {
		if phone == conversation.PhoneDestination {
			return true
		}
	}

	return false
}

// IsUserOnConversationByID - Indicates if the user is on the conversation, fetch conversation first.
func IsUserOnConversationByID(table dynamo.Table, conversationID string, user *User, debug bool) (bool, error) {
	conversation, err := GetConversationByID(table, conversationID, debug)
	if err != nil {
		return false, err
	}

	return IsUserOnConversation(conversation, user, debug), nil
}

// GetUsersOnConversation - Get the allowed users that can participate in the conversation from the company side
func GetUsersOnConversation(cognitoID string, phoneDestination string) ([]string, error) {
	phoneRe := regexp.MustCompile(`^\+(\d+)$`)
	phoneMatch := phoneRe.FindStringSubmatch(phoneDestination)
	phone := phoneDestination
	if len(phoneMatch) == 2 {
		phone = phoneMatch[1]
	}
	groupName := fmt.Sprintf("group_%s", phone)
	return GetUsersInGroup(cognitoID, groupName)
}

// GetUsersInGroup - Get users from group
func GetUsersInGroup(cognitoID string, groupName string) ([]string, error) {
	var userIDs []string
	cognito := cognitoidentityprovider.New(session.Must(session.NewSession()))
	users, err := cognito.ListUsersInGroup(&cognitoidentityprovider.ListUsersInGroupInput{
		UserPoolId: &cognitoID,
		GroupName:  &groupName,
	})

	if err != nil {
		if err == dynamo.ErrNotFound {
			return userIDs, nil
		}
		return nil, err
	}

	for _, userGroup := range users.Users {
		if *userGroup.UserStatus == "CONFIRMED" {
			userIDs = append(userIDs, *userGroup.Username)
		}
	}
	return userIDs, nil
}

// SetBoostConversationID - Set the Boost conversation id
func SetBoostConversationID(table dynamo.Table, conversationID string, boostConversationID string, debug bool) error {
	conversation, err := GetConversationByID(table, conversationID, debug)
	if err != nil {
		log.Println("[SetBoostConversationID] Error while getting conversation by id")
		return err
	}
	conversation.BoostConversationID = boostConversationID
	_, saveErr := SaveConversation(table, conversation)
	return saveErr
}

// AddConversation - Add conversation from guest.
//   Right now assume from guest, so always setting the DestinationGUID as companyID.
func AddConversation(table dynamo.Table, message *EventMessage, companyID string) (*Conversation, error) {
	// TODO: move away from guid and have these values be computed
	// ConversationID: PhoneGuest + PhoneDestination
	// Drop GuestGuid and DestinationGuid?
	conversation := Conversation{
		PhoneGuest:           message.FromPhoneNumber,
		PhoneDestination:     message.ToPhoneNumber,
		ConversationID:       uuid.New().String(),
		GuestGUID:            uuid.New().String(),
		DestinationGUID:      companyID,
		CreatedAt:            utils.Now(),
		BotEnabled:           true,
		BoostConversationID:  "",
		GuestName:            "",
		GuestNotes:           "",
		CheckedAt:            "",
		Read:                 false,
		LastMessage:          message.Body,
		LastMessageCreatedAt: message.CreatedAt,
	}
	return SaveConversation(table, &conversation)
}

// SaveConversation - save (create/update) conversation
func SaveConversation(table dynamo.Table, conversation *Conversation) (*Conversation, error) {
	return conversation, table.Put(conversation).Run()
}
