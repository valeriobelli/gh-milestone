package gh

import (
	"os"
	"strings"
)

func RetrieveCurrentToken() string {
	var getHost = func() string {
		env_host := os.Getenv("GITHUB_MILESTONE_HOST")

		if env_host == "" {
			return "github.com"
		}

		return env_host
	}

	token, err := Execute([]string{"config", "get", "-h", getHost(), "oauth_token"})

	if err != nil {
		return ""
	}

	return strings.TrimRight(token, "\n")
}
