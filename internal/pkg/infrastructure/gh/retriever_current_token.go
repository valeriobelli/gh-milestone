package gh

import (
	"os"

	"github.com/cli/go-gh/v2/pkg/auth"
)

func RetrieveCurrentToken() string {
	var getHost = func() string {
		env_host := os.Getenv("GITHUB_MILESTONE_HOST")

		if env_host == "" {
			return "github.com"
		}

		return env_host
	}

	token, _ := auth.TokenForHost(getHost())

	return token
}
