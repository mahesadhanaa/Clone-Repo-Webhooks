package lambdaChitato

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"gopkg.in/go-playground/webhooks.v5/github"
	"gopkg.in/src-d/go-git.v4"
	gitHTTP "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

// router function
func Router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return create(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

// function for POST method
func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	pullRequestPayload := new(github.PullRequestPayload)
	issueCommentPayload := new(github.IssueCommentPayload)

	err := json.Unmarshal([]byte(req.Body), pullRequestPayload)

	username := os.Getenv("username")
	password := os.Getenv("password")

	//status pull request opened
	if pullRequestPayload.Action == "opened" {
		// clone repository
		CloneRepository("/tmp/chitato-repository/"+pullRequestPayload.Repository.Name, pullRequestPayload.Repository.HTMLURL, username, password)

		// whatever do something

		// status pull request closed
	} else if pullRequestPayload.Action == "closed" {
		// clone repository
		CloneRepository("/tmp/chitato-repository/"+pullRequestPayload.Repository.Name, pullRequestPayload.Repository.HTMLURL, username, password)

		// whatever do something

	} else if issueCommentPayload.Action == "created" && issueCommentPayload.Comment.Body == "chitato check" {
		// clone repository
		CloneRepository("/tmp/chitato-repository/"+pullRequestPayload.Repository.Name, pullRequestPayload.Repository.HTMLURL, username, password)

		// whatever do something
	}

	if err != nil {
		return clientError(http.StatusUnprocessableEntity)

	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "success for Action " + pullRequestPayload.Action,
		Headers:    map[string]string{"Location": fmt.Sprintf("/webhooks?action=%s", pullRequestPayload.Action)},
	}, nil
}

// server error function
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// client error function
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

// funcion clone repository
func CloneRepository(pathDir string, url string, username string, password string) {
	_, err := git.PlainClone(pathDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
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
