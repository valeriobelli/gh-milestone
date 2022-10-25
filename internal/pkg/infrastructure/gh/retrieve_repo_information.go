package gh

import (
	"errors"
	"strings"
)

type RepoInfo struct {
	Name  string
	Owner string
}

func RetrieveRepoInformation() (*RepoInfo, error) {
	owner, err := Execute([]string{"repo", "view", "--json", "owner", "--jq", ".owner.login"})

	if err != nil {
		return nil, errors.New("failed to retrieve the owner of the repository")
	}

	name, err := Execute([]string{"repo", "view", "--json", "name", "--jq", ".name"})

	if err != nil {
		return nil, errors.New("failed to retrieve the name of the repository")
	}

	return &RepoInfo{Name: strings.TrimSpace(name), Owner: strings.TrimSpace(owner)}, nil
}
