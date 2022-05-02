package create

import (
	"context"
	"fmt"
	"time"

	ghub "github.com/google/go-github/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/gh"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/http"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/spinner"
)

type CreateMilestoneConfig struct {
	DueDate     *time.Time
	Title       *string
	Description *string
	Verbose     bool
}

type CreateMilestone struct {
	config CreateMilestoneConfig
}

func NewCreateMilestone(config CreateMilestoneConfig) *CreateMilestone {
	return &CreateMilestone{config: config}
}

func (cm CreateMilestone) Execute() {
	repoInfo, err := gh.RetrieveRepoInformation()

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	client := github.NewRestClient(http.NewClient())

	spinner := spinner.NewSpinner()

	spinner.Start()

	_, _, err = client.Issues.CreateMilestone(
		context.Background(),
		repoInfo.Owner,
		repoInfo.Name,
		&ghub.Milestone{
			Description: cm.config.Description,
			DueOn:       cm.config.DueDate,
			Title:       cm.config.Title,
		},
	)

	spinner.Stop()

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	fmt.Println("Milestone has been created.")
}
