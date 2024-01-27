package utils

import "github.com/aws/aws-sdk-go/service/dynamodb"

// DereferenceString - dereference string in a nil safe way
func DereferenceString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// DereferenceBool - dereference bool in a nil safe way
func DereferenceBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

// ExtractDynamoResultString - safely extract the dynamo result string
func ExtractDynamoResultString(item map[string]*dynamodb.AttributeValue, key string) string {
	if item == nil {
		return ""
	}

	value := item[key]
	if value == nil {
		return ""
	}

	return DereferenceString(value.S)
}

// ExtractDynamoResultBool - safely extract the dynamo result bool
func ExtractDynamoResultBool(item map[string]*dynamodb.AttributeValue, key string) bool {
	if item == nil {
		return false
	}

	value := item[key]
	if value == nil {
		return false
	}

	return DereferenceBool(value.BOOL)
}
