package gh

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type RepoInfo struct {
	Name  string
	Owner string
}

func extractRepoInformationFromRepoString(repoString string) (*RepoInfo, error) {
	if repoString == "" {
		return nil, nil
	}

	pattern, err := regexp.Compile("(?:(?:https://www.|https://|www.)?github.com/)?([a-zA-Z0-9_-]+)/([a-zA-Z0-9_-]+)")

	if err != nil {
		return nil, err
	}

	subStr := pattern.FindStringSubmatch(repoString)

	if len(subStr) == 0 {
		return nil, fmt.Errorf("the repo option \"%s\" does not respect the format [HOST/]OWNER/REPO", repoString)
	}

	fmt.Printf("%v", subStr)

	return &RepoInfo{
		Name:  subStr[2],
		Owner: subStr[1],
	}, nil
}

func RetrieveRepoInformation(repoString string) (*RepoInfo, error) {
	repoInfoFromRepoString, err := extractRepoInformationFromRepoString(repoString)

	if err != nil {
		return nil, err
	}

	if repoInfoFromRepoString != nil {
		return repoInfoFromRepoString, nil
	}

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
