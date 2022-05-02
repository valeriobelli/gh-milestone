package edit

import (
	"context"
	"fmt"
	"time"

	ghub "github.com/google/go-github/v44/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/gh"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/http"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/spinner"
)

type EditMilestoneConfig struct {
	Description *string
	DueDate     *time.Time
	State       *string
	Title       *string
}

type EditMilestone struct {
	config EditMilestoneConfig
}

func NewEditMilestone(config EditMilestoneConfig) *EditMilestone {
	return &EditMilestone{config: config}
}

func (em EditMilestone) Execute(number int) {
	repoInfo, err := gh.RetrieveRepoInformation()

	if err != nil {
		fmt.Println(err.Error())

		return
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
			DueOn:       em.config.DueDate,
			State:       em.config.State,
			Title:       em.config.Title,
		},
	)

	spinner.Stop()

	if err != nil {
		fmt.Println(handleResponseError(response).Error())

		return
	}

	fmt.Println(*milestone.HTMLURL)
}
