package alexa

// Response is the alexa response
type Response struct {
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          ResponseBody           `json:"response"`
	Version           string                 `json:"version"`
}

// ResponseBody is the alexa response body
type ResponseBody struct {
	Card             *ResponsePayload `json:"card,omitempty"`
	OutputSpeech     *ResponsePayload `json:"outputSpeech,omitempty"`
	Reprompt         *Reprompt        `json:"reprompt,omitempty"`
	ShouldEndSession bool             `json:"shouldEndSession"`
}

// ResponsePayload is the response payload
type ResponsePayload struct {
	Content string        `json:"content,omitempty"`
	Image   ResponseImage `json:"image,omitempty"`
	Text    string        `json:"text,omitempty"`
	Title   string        `json:"title,omitempty"`
	Type    string        `json:"type,omitempty"`
	SSML    string        `json:"ssml,omitempty"`
}

// Reprompt is the response for a reprompt
type Reprompt struct {
	OutputSpeech ResponsePayload `json:"outputSpeech,omitempty"`
}

// ResponseImage is a response image
type ResponseImage struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

// NewResponse returns a new Alexa response
func NewResponse() *Response {
	return &Response{
		Version: "1.0",
		Response: ResponseBody{
			ShouldEndSession: true,
		},
	}
}

// OutputSpeech sets the response output speech to plain text
func (r *Response) OutputSpeech(text string) *Response {
	r.Response.OutputSpeech = &ResponsePayload{
		Type: "PlainText",
		Text: text,
	}

	return r
}

// OutputSpeechSSML sets the response output speech to SSML
func (r *Response) OutputSpeechSSML(text string) *Response {
	r.Response.OutputSpeech = &ResponsePayload{
		Type: "SSML",
		SSML: text,
	}

	return r
}

// Card sets the response card
func (r *Response) Card(title string, content string) *Response {
	return r.SimpleCard(title, content)
}

// SimpleCard sets the response card to a simple card
func (r *Response) SimpleCard(title string, content string) *Response {
	r.Response.Card = &ResponsePayload{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}

	return r
}

// StandardCard sets the response card to a standard card
func (r *Response) StandardCard(title string, content string, smallImg string, largeImg string) *Response {
	r.Response.Card = &ResponsePayload{
		Type:    "Standard",
		Title:   title,
		Content: content,
	}

	if smallImg != "" {
		r.Response.Card.Image.SmallImageURL = smallImg
	}

	if largeImg != "" {
		r.Response.Card.Image.LargeImageURL = largeImg
	}

	return r
}

// LinkAccountCard send a link account card
func (r *Response) LinkAccountCard() *Response {
	r.Response.Card = &ResponsePayload{
		Type: "LinkAccount",
	}

	return r
}

// Reprompt returns a reprompt response
func (r *Response) Reprompt(text string) *Response {
	r.Response.Reprompt = &Reprompt{
		OutputSpeech: ResponsePayload{
			Type: "PlainText",
			Text: text,
		},
	}

	return r
}

// EndSession ends the session
func (r *Response) EndSession(flag bool) *Response {
	r.Response.ShouldEndSession = flag

	return r
}
