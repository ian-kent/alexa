package main

import (
	"fmt"
	"os"

	"github.com/apex/go-apex"
	"github.com/ian-kent/alexa"
)

func main() {
	app := alexa.Application{
		ID: os.Getenv("ALEXA_APP_ID"),
		LaunchHandler: func(req *alexa.Request, res *alexa.Response) {
			res.OutputSpeech("Go app test")
		},
		IntentHandler: func(req *alexa.Request, res *alexa.Response) {
			switch req.IntentName() {
			case "Ask":
				v, err := req.SlotValue("topic")
				if err != nil {
					res.OutputSpeech("Internal error")
					return
				}
				if v == "" {
					res.OutputSpeech("I'm not sure what you mean by that")
					return
				}

				switch v {
				case "go":
					res.OutputSpeech(fmt.Sprintf("Good choice, that's all you need to know!"))
				default:
					res.OutputSpeech(fmt.Sprintf("Don't use %s, use go!", v))
				}
			default:
				res.OutputSpeech("Unexpected request")
			}
		},
		SessionEndedHandler: func(req *alexa.Request, res *alexa.Response) {
		},
	}
	apex.HandleFunc(app.Handler)
}
