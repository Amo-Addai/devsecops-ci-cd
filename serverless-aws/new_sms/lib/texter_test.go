package lib

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"net/http"
	"net/url"
	"testing"
)

func TestNewTexter(t *testing.T) {
	phone := "+15555551111"
	authToken := "34485e66-fc0d-4a7a-bc08-4b0946fa0218"
	accountSid := "3dfdf76d-4da0-4f9b-b967-2a0801ceefb8"
	thisURL := "http://example.com"
	data := url.Values{}
	data.Add("testing", "yes")
	debug := false
	texter := NewTexter(phone, authToken, accountSid, thisURL, data, debug)

	assert.NotNil(t, texter.Twilio)
	assert.Equal(t, phone, texter.TwilioPhoneNumber)
	assert.Equal(t, authToken, texter.TwilioAuthToken)
	assert.Equal(t, accountSid, texter.TwilioAccountSid)
	assert.Equal(t, thisURL, texter.ThisURL)
	assert.Equal(t, data, texter.Data)
	assert.Equal(t, debug, texter.Debug)
}

func TestValidateWebhook(t *testing.T) {
	phone := "+15555551111"
	authToken := "34485e66-fc0d-4a7a-bc08-4b0946fa0218"
	accountSid := "3dfdf76d-4da0-4f9b-b967-2a0801ceefb8"
	thisURL := "http://example.com"
	data := url.Values{}
	debug := false
	texter := NewTexter(phone, authToken, accountSid, thisURL, data, debug)

	mock := &TwilioMock{}
	texter.Twilio = mock
	request := events.APIGatewayProxyRequest{}

	req := http.Request{}
	baseURL, parseErr := url.Parse(thisURL)
	assert.Nil(t, parseErr)
	req.URL = baseURL
	req.Method = "POST"
	req.Form = data
	req.PostForm = data
	req.ParseForm()
	req.Header = make(map[string][]string)

	mock.On("CheckRequestSignature", &req, "").Return(true, nil)
	valid := texter.ValidateWebhook(request)
	assert.True(t, valid)
	mock.AssertExpectations(t)
}

type TwilioMock struct {
	mock.Mock
}

func (m *TwilioMock) CheckRequestSignature(req *http.Request, baseURL string) (bool, error) {
	m.Called(req, baseURL)
	return true, nil
}
