package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
)

var r = regexp.MustCompile(`@[a-zA-Z0-9_\-]+`)

type Webhook struct {
	EventType       string
	Payload         []byte
	Repository      string
	Title           string
	OriginComment   string
	ReplacedComment string
	User            string
	HTMLURL         string
}

var secretGithub = os.Getenv("GITHUB_SECRET")

func GetGithubComment(webhook *Webhook) {
	switch webhook.EventType {
	case "issue_comment":
		IssueComment(webhook)
	case "pull_request_review_comment":
		PullRequestComment(webhook)
	default:
		panic("Event doesn't exist")
	}

}

func IssueComment(webhook *Webhook) {
	var issueGithub github.IssueCommentEvent
	err := json.Unmarshal(webhook.Payload, &issueGithub)
	if err != nil {
		panic(err)
	}

	webhook.Repository = *issueGithub.Repo.Name
	webhook.Title = *issueGithub.Issue.Title
	webhook.User = *issueGithub.Comment.User.Login
	webhook.OriginComment = *issueGithub.Comment.Body
	webhook.HTMLURL = *issueGithub.Comment.HTMLURL
}

func PullRequestComment(webhook *Webhook) {
	var pullrequestGithub github.PullRequestReviewCommentEvent
	err := json.Unmarshal(webhook.Payload, &pullrequestGithub)
	if err != nil {
		panic(err)
	}

	webhook.Repository = *pullrequestGithub.Repo.Name
	webhook.User = *pullrequestGithub.Comment.User.Login
	webhook.Title = *pullrequestGithub.PullRequest.Title
	webhook.OriginComment = *pullrequestGithub.Comment.Body
	webhook.HTMLURL = *pullrequestGithub.Comment.HTMLURL
}

// ReplaceComment replace github account to slack
func ReplaceComment(comment string, conf *Config) string {

	matches := r.FindAllStringSubmatch(comment, -1)
	for _, val := range matches {
		slackName, _ := conf.Accounts[val[0]]
		comment = strings.Replace(comment, val[0], slackName, -1)
	}
	return comment
}
