package gh

import (
	gogh "github.com/cli/go-gh"
)

func Execute(args []string) (string, error) {
	result, _, err := gogh.Exec(args...)

	return result.String(), err
}
