package lib

import (
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sfreiberg/gotwilio"
)

// TextingService - interface to make texts
type TextingService interface {
	CheckRequestSignature(*http.Request, string) (bool, error)
}

// Texter - Object in charge of texting
type Texter struct {
	Twilio            TextingService
	TwilioPhoneNumber string
	TwilioAuthToken   string
	TwilioAccountSid  string
	ThisURL           string
	Data              url.Values
	Debug             bool
}

// NewTexter Creates a new texting object
func NewTexter(phone string, authToken string, accountSid string, thisURL string, data url.Values, debug bool) Texter {
	return Texter{
		gotwilio.NewTwilioClient(accountSid, authToken),
		phone,
		authToken,
		accountSid,
		thisURL,
		data,
		debug}
}

// ValidateWebhook - Validate the webhook call
func (texter *Texter) ValidateWebhook(request events.APIGatewayProxyRequest) bool {
	httpRequest := createHTTPRequest(request, texter.ThisURL, texter.Data)
	valid, err := texter.Twilio.CheckRequestSignature(&httpRequest, "")

	if texter.Debug {
		log.Printf("%+v\n", httpRequest)
		log.Printf("%+v\n", httpRequest.Form)
		log.Printf("%+v\n", httpRequest.PostForm)
		log.Printf("webhook valid: %t", valid)
		log.Print("error", err)
	}
	if err != nil {
		valid = false
	}

	return valid
}

func createHTTPRequest(request events.APIGatewayProxyRequest, thisURL string, data url.Values) http.Request {
	r := http.Request{}
	baseURL, err := url.Parse(thisURL)
	if err != nil {
		log.Println("Malformed URL: ", err.Error())
	}
	r.URL = baseURL

	r.Header = make(map[string][]string)
	for k, v := range request.Headers {
		r.Header.Set(k, v)
	}
	r.Form = data
	r.PostForm = data
	r.ParseForm()
	r.Method = "POST"
	return r
}
