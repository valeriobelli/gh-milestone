package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/gh"
	"golang.org/x/oauth2"
)

func NewClient() *http.Client {
	var getToken = func() string {
		github_token := strings.TrimRight(os.Getenv("GITHUB_TOKEN"), "\n")

		if github_token == "" {
			return gh.RetrieveCurrentToken()
		}

		return github_token
	}

	token := getToken()

	fmt.Println()

	token_source := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	return oauth2.NewClient(context.Background(), token_source)
}
