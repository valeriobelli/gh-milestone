package create

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

type CreateMilestoneConfig struct {
	Description string
	DueDate     *time.Time
	Title       string
}

type CreateMilestone struct {
	config CreateMilestoneConfig
}

func NewCreateMilestone(config CreateMilestoneConfig) *CreateMilestone {
	return &CreateMilestone{config: config}
}

func (cm CreateMilestone) Execute() error {
	repoInfo, err := gh.RetrieveRepoInformation()

	if err != nil {
		return err
	}

	fmt.Printf("Creating milestone in %s/%s\n\n", repoInfo.Owner, repoInfo.Name)

	survey := NewSurvey(Flags{
		Description: cm.config.Description,
		DueDate:     cm.config.DueDate,
		Title:       cm.config.Title,
	})

	answers, err := survey.Ask()

	if err != nil {
		return err
	}

	if !answers.Confirm {
		fmt.Println("Discarding.")

		return nil
	}

	client := github.NewRestClient(http.NewClient())

	spinner := spinner.NewSpinner()

	spinner.Start()

	milestone, response, err := client.Issues.CreateMilestone(
		context.Background(),
		repoInfo.Owner,
		repoInfo.Name,
		&ghub.Milestone{
			Description: &answers.Description,
			DueOn:       answers.getTime(),
			Title:       &answers.Title,
		},
	)

	spinner.Stop()

	if err != nil {
		if response == nil {
			return err
		}

		return handleResponseError(response)
	}

	fmt.Println(milestone.GetHTMLURL())

	return nil
}
