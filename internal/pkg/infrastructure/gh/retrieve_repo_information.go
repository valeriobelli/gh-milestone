package gh

import "strings"

type RepoInfo struct {
	Name  string
	Owner string
}

func RetrieveRepoInformation() (RepoInfo, error) {
	owner, err := Execute([]string{"repo", "view", "--json", "owner", "--jq", ".owner.login"})
	name, err := Execute([]string{"repo", "view", "--json", "name", "--jq", ".name"})

	return RepoInfo{Name: strings.TrimSpace(name), Owner: strings.TrimSpace(owner)}, err
}
