package gh

import (
	"github.com/cli/go-gh"
)

func Execute(args []string) (string, error) {
	result, _, err := gh.Exec(args...)

	return result.String(), err
}
