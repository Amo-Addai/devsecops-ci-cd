package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// User - pulled from user claims
type User struct {
	CognitoGroups     []string `json:"cognito_groups"`
	DestinationPhones []string `json:"destination_phones"`
	Superadmin        bool
	CognitoID         string
	ID                string
}

// ParseUser - parse user from authorizer
func ParseUser(authorizer map[string]interface{}) (*User, error) {
	if authorizer["claims"] == nil {
		return nil, errors.New("Missing claim")
	}
	claims := authorizer["claims"].(map[string]interface{})

	if claims["iss"] == nil || claims["sub"] == nil {
		return nil, errors.New("Missing claim")
	}
	cognitoURL := claims["iss"].(string)
	id := claims["sub"].(string)

	cognitoRe := regexp.MustCompile(`^https\:\/\/cognito-idp\.[\w\-]+\.amazonaws\.com\/([\w\-\_]+)$`)
	cognitoMatch := cognitoRe.FindStringSubmatch(cognitoURL)
	cognitoID := ""
	if len(cognitoMatch) == 2 {
		cognitoID = cognitoMatch[1]
	}

	cognitoGroups := claims["cognito:groups"].(string)
	groups := strings.Split(cognitoGroups, ",")
	var destinationPhones []string
	superadmin := false

	re := regexp.MustCompile(`^group_(\d+)$`)
	for _, group := range groups {
		if group == "superadmin" {
			superadmin = true
		} else {
			match := re.FindStringSubmatch(group)
			if len(match) == 2 {
				destinationPhones = append(destinationPhones, fmt.Sprintf("+%s", match[1]))
			}
		}
	}

	return &User{
		CognitoGroups:     groups,
		DestinationPhones: destinationPhones,
		Superadmin:        superadmin,
		CognitoID:         cognitoID,
		ID:                id,
	}, nil
}
