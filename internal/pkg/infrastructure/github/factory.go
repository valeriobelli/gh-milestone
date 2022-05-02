package github

import (
	"net/http"

	"github.com/google/go-github/v44/github"
	"github.com/shurcooL/githubv4"
)

func NewGraphQlClient(httpClient *http.Client) *githubv4.Client {
	return githubv4.NewClient(httpClient)
}

func NewRestClient(httpClient *http.Client) *github.Client {
	return github.NewClient(httpClient)
}
