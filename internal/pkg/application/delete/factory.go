package delete

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/gh"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/http"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/spinner"
)

type DeleteMilestoneConfig struct {
	Confirm bool
	Repo    string
}

type DeleteMilestone struct {
	config DeleteMilestoneConfig
}

func NewDeleteMilestone(config DeleteMilestoneConfig) *DeleteMilestone {
	return &DeleteMilestone{config: config}
}

func (em DeleteMilestone) Execute(number int) error {
	repoInfo, err := gh.RetrieveRepoInformation(em.config.Repo)

	if err != nil {
		return err
	}

	client := github.NewRestClient(http.NewClient())

	spinner := spinner.NewSpinner()

	spinner.Start()

	milestone, response, err := client.Issues.GetMilestone(context.Background(), repoInfo.Owner, repoInfo.Name, number)

	spinner.Stop()

	if err != nil {
		return handleResponseError(response)
	}

	survey := NewSurvey(Config{
		Milestone: milestone,
		Confirm:   em.config.Confirm,
	})

	surveyAnswers, err := survey.Ask()

	if err != nil {
		return err
	}

	if !surveyAnswers.Confirm {
		fmt.Printf("Milestone #%d was not deleted.\n", *milestone.Number)

		return nil
	}

	spinner.Start()

	response, err = client.Issues.DeleteMilestone(
		context.Background(),
		repoInfo.Owner,
		repoInfo.Name,
		number,
	)

	spinner.Stop()

	if err != nil {
		return handleResponseError(response)
	}

	fmt.Printf(color.RedString("✔ ")+"Deleted milestone #%d (%s).\n", *milestone.Number, *milestone.Title)

	return nil
}
