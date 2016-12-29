package alexa

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"github.com/apex/go-apex"
)

// Application is an alexa application
type Application struct {
	ID                  string
	LaunchHandler       func(*Request, *Response)
	IntentHandler       func(*Request, *Response)
	SessionEndedHandler func(*Request, *Response)
}

// Handler is a http.Handler for the application
func (a Application) Handler(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
	log.Println(string(event))
	req, err := ParseRequest(a.ID, bytes.NewReader(event))
	if err != nil {
		return nil, err
	}

	res := NewResponse()

	switch req.Type() {
	case "LaunchRequest":
		if a.LaunchHandler != nil {
			a.LaunchHandler(req, res)
		}
	case "IntentRequest":
		if a.IntentHandler != nil {
			a.IntentHandler(req, res)
		}
	case "SessionEndedRequest":
		if a.SessionEndedHandler != nil {
			a.SessionEndedHandler(req, res)
		}
	default:
		return nil, errors.New("invalid request type")
	}

	return res, nil
}
