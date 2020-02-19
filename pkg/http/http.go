package http

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/go-playground/webhooks.v5/github"
	"gopkg.in/src-d/go-git.v4"
	gitHTTP "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func WebhookHandler(r *http.Request) (interface{}, error) {

	secret := os.Getenv("SECRET")
	hook, _ := github.New(github.Options.Secret(secret))
	payload, err := hook.Parse(r, github.PullRequestEvent, github.IssueCommentEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn't one of the ones asked to be parsed
		}
	}

	switch payload.(type) {

	case github.IssueCommentPayload:
		issueComment := payload.(github.IssueCommentPayload)

		if issueComment.Action == "created" && issueComment.Comment.Body == "chitato check" {

			// whatever do something
		}

	case github.PullRequestPayload:
		pullRequest := payload.(github.PullRequestPayload)
		username := os.Getenv("username")
		password := os.Getenv("password")

		// Action condition opened
		if pullRequest.Action == "opened" {

			CloneRepository("/tmp/chitato-repository/"+pullRequest.Repository.Name, pullRequest.Repository.HTMLURL, username, password)
			// whatever do something

			// Action condition closed
		} else if pullRequest.Action == "closed" {

			CloneRepository("/tmp/chitato-repository/"+pullRequest.Repository.Name, pullRequest.Repository.HTMLURL, username, password)

			// whatever do something
		}

	}
	return nil, nil
}

// function clone repository
func CloneRepository(pathDir string, url string, username string, password string) {
	_, err := git.PlainClone(pathDir, false, &git.CloneOptions{
		URL: url,
		Auth: &gitHTTP.BasicAuth{
			Username: username,
			Password: password,
		},
	})
	fmt.Println(err)

	if err == git.ErrRepositoryAlreadyExists {
		os.RemoveAll(pathDir)
		_, err = git.PlainClone(pathDir, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
			Auth: &gitHTTP.BasicAuth{
				Username: username,
				Password: password,
			},
		})
	}
}
