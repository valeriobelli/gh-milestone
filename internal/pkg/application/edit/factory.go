package edit

import (
	"context"
	"fmt"
	"time"

	ghub "github.com/google/go-github/v68/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/gh"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/http"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/spinner"
)

type EditMilestoneConfig struct {
	Description *string
	DueDate     *time.Time
	Repo        string
	State       *string
	Title       *string
}

type EditMilestone struct {
	config EditMilestoneConfig
}

func NewEditMilestone(config EditMilestoneConfig) *EditMilestone {
	return &EditMilestone{config: config}
}

func (em EditMilestone) Execute(number int) error {
	repoInfo, err := gh.RetrieveRepoInformation(em.config.Repo)

	if err != nil {
		return err
	}

	client := github.NewRestClient(http.NewClient())

	spinner := spinner.NewSpinner()

	spinner.Start()

	milestone, response, err := client.Issues.EditMilestone(
		context.Background(),
		repoInfo.Owner,
		repoInfo.Name,
		number,
		&ghub.Milestone{
			Description: em.config.Description,
			DueOn:       toTimestamp(em.config.DueDate),
			State:       em.config.State,
			Title:       em.config.Title,
		},
	)

	spinner.Stop()

	if err != nil {
		return handleResponseError(response)
	}

	fmt.Println(*milestone.HTMLURL)

	return nil
}

func toTimestamp(t *time.Time) *ghub.Timestamp {
	if t == nil {
		return nil
	}

	return &ghub.Timestamp{Time: *t}
}
