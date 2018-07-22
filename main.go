package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
)

func init() {
	http.HandleFunc("/github/events", GitEventHandler)
}

func main() {
	var port string
	if os.Getenv("PORT") == "" {
		port = "3000"
	} else {
		port = os.Getenv("PORT")
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func GitEventHandler(w http.ResponseWriter, r *http.Request) {

	webhook := Webhook{}
	webhook.EventType = r.Header.Get("X-GitHub-Event")
	payload, err := github.ValidatePayload(r, []byte(secretGithub))

	if err != nil {
		panic("Invalid signature")
	}

	webhook.Payload = payload

	GetGithubComment(&webhook)

	conf, err := ParseFile("./config.json")
	if err != nil {
		panic("Invalid Config")
	}

	comment := ReplaceComment(webhook.OriginComment, conf)

	if comment == webhook.OriginComment {
		return
	}

	var text string
	text = fmt.Sprintf("%v *【%v】%v* \n", text, webhook.Repository, webhook.Title)
	text = fmt.Sprintf("%v%v\n", text, webhook.HTMLURL)
	text = fmt.Sprintf("%v>Comment created by: %v\n", text, webhook.User)
	text = fmt.Sprintf("%v\n%v\n", text, comment)

	sendToSlack(text)
}
