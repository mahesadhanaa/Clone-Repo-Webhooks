package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	lambdaChitato "github.com/mahesadhanaa/go-git-lambdaAndHttp/pkg/lambda"
)

func main() {
	lambda.Start(lambdaChitato.Router)
}
