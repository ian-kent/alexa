package alexa

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// ErrSlotNotFound is returned when a slot is not found
var ErrSlotNotFound = errors.New("slot not found")

// ErrRequestExpired is returned when the request is >150s old
var ErrRequestExpired = errors.New("request timestamp has expired")

// ErrInvalidAppID is returned when the application ID is invalid
var ErrInvalidAppID = errors.New("invalid application ID")

// Request is an alexa request
type Request struct {
	Version string      `json:"version"`
	Session Session     `json:"session"`
	Request RequestBody `json:"request"`
}

// Session is an alexa session
type Session struct {
	New         bool   `json:"new"`
	SessionID   string `json:"sessionId"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	Attributes struct {
		String map[string]interface{} `json:"string"`
	} `json:"attributes"`
	User struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"user"`
}

// RequestBody is the request body
type RequestBody struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Intent    Intent `json:"intent,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

// Intent is the request intent
type Intent struct {
	Name  string          `json:"name"`
	Slots map[string]Slot `json:"slots"`
}

// Slot is a slot value
type Slot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ParseRequest parses JSON bytes into a *Request an verifies it
func ParseRequest(appID string, r io.Reader) (req *Request, err error) {
	err = json.NewDecoder(r).Decode(&req)
	if err != nil {
		return
	}

	if !req.VerifyTimestamp() {
		err = ErrRequestExpired
		return
	}

	if !req.VerifyApplicationID(appID) {
		err = ErrInvalidAppID
		return
	}

	return
}

// VerifyTimestamp verifies the request timestamp
func (r *Request) VerifyTimestamp() bool {
	reqTimestamp, err := time.Parse("2006-01-02T15:04:05Z", r.Request.Timestamp)
	if err != nil {
		return false
	}

	if time.Since(reqTimestamp) > time.Duration(150)*time.Second {
		return false
	}

	return true
}

// VerifyApplicationID verifies the application ID
func (r *Request) VerifyApplicationID(appID string) bool {
	// this appears to be optional in the nodejs library
	if len(appID) == 0 {
		return true
	}

	if r.Session.Application.ApplicationID == appID {
		return true
	}

	return false
}

// SessionID gets the session ID
func (r *Request) SessionID() string {
	return r.Session.SessionID
}

// UserID gets the user ID
func (r *Request) UserID() string {
	return r.Session.User.UserID
}

// Type gets the request type
func (r *Request) Type() string {
	return r.Request.Type
}

// IntentName gets the intent name
func (r *Request) IntentName() string {
	if r.Type() == "IntentRequest" {
		return r.Request.Intent.Name
	}

	return r.Type()
}

// SlotValue returns a slot value
func (r *Request) SlotValue(slotName string) (string, error) {
	if v, ok := r.Request.Intent.Slots[slotName]; ok {
		return v.Value, nil
	}

	return "", ErrSlotNotFound
}

// Slots returns all slots
func (r *Request) Slots() map[string]Slot {
	return r.Request.Intent.Slots
}
