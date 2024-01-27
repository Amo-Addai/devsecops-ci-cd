package lib

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestGetEnvironment(t *testing.T) {
	environment := NewEnvironment()

	assert.Equal(t, false, environment.Debug)
	assert.Equal(t, "", environment.TwilioPhoneNumber)
	assert.Equal(t, "", environment.TwilioAccountSid)
	assert.Equal(t, "", environment.ThisURL)
	assert.Equal(t, "", environment.TopicURL)
}
