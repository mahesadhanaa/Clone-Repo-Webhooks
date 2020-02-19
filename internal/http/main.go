package main

import (
	"net/http"

	server "github.com/mahesadhana/go-git-lambdaAndHttp/pkg/http"
)

func main() {
	http.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		server.http.WebhookHandler(r)
	})
	http.ListenAndServe(":6969", nil)
}