package http

import (
	"context"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

func NewClient() *http.Client {
	token_source := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	return oauth2.NewClient(context.Background(), token_source)
}
