package view

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/shurcooL/githubv4"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	github_entities "github.com/valeriobelli/gh-milestone/internal/pkg/domain/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/gh"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/http"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/spinner"
)

var query struct {
	Repository struct {
		*github_entities.Milestone `graphql:"milestone(number: $number)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

type ViewMilestoneConfig struct {
	Web bool
}

type ViewMilestone struct {
	config ViewMilestoneConfig
}

func NewViewMilestone(config ViewMilestoneConfig) *ViewMilestone {
	return &ViewMilestone{config: config}
}

func (vm ViewMilestone) Execute(number int) {
	repoInfo, err := gh.RetrieveRepoInformation()

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	client := github.NewGraphQlClient(http.NewClient())

	spinner := spinner.NewSpinner()

	spinner.Start()

	err = client.Query(context.Background(), &query, map[string]interface{}{
		"name":   githubv4.String(strings.TrimSpace(repoInfo.Name)),
		"number": githubv4.Int(number),
		"owner":  githubv4.String(strings.TrimSpace(repoInfo.Owner)),
	})

	spinner.Stop()

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	milestone := query.Repository.Milestone

	if milestone == nil {
		fmt.Println("Requested milestone does not exist.")

		return
	}

	if vm.config.Web {
		vm.openUrl(*milestone)

		return
	}

	vm.printConsole(*milestone)
}

func (vm ViewMilestone) mapTitle(milestone github_entities.Milestone) string {
	if milestone.Title == "" {
		return "<No Title>"
	}

	return milestone.Title
}

func (vm ViewMilestone) mapDescription(milestone github_entities.Milestone) string {
	if milestone.Description == "" {
		return "<No description>"
	}

	return milestone.Description
}

func (vm ViewMilestone) mapDueOn(milestone github_entities.Milestone) string {
	parsedTime, err := time.Parse(milestone.DueOn, constants.DateFormat)

	if err != nil {
		return "<No due date>"
	}

	return parsedTime.Format(constants.DateFormat)
}

func (vm ViewMilestone) openUrl(milestone github_entities.Milestone) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	args = append(args, milestone.Url)

	exec.Command(cmd, args...).Start()
}

func (vm ViewMilestone) printConsole(milestone github_entities.Milestone) {
	color.Set(color.FgHiWhite)

	fmt.Print(milestone.Title)
	fmt.Print(" ")

	if milestone.Closed {
		color.Set(color.FgRed)
	} else {
		color.Set(color.FgGreen)
	}

	fmt.Print(milestone.State)

	color.Unset()

	fmt.Print(" - ")

	fmt.Printf("%d%% complete\n", int(milestone.ProgressPercentage))

	color.Set(color.FgWhite)

	if milestone.DueOn != "" {
		parsedTime, err := time.Parse(time.RFC3339, milestone.DueOn)

		if err != nil {
			fmt.Print(err)
			return
		}

		year, month, day := parsedTime.Date()

		fmt.Printf("Due by %d-%d-%d - ", year, month, day)
	} else {
		fmt.Print("No due date - ")
	}

	parsedTime, err := time.Parse(time.RFC3339, milestone.UpdatedAt)

	if err != nil {
		fmt.Print(err)
		return
	}

	year, month, day := parsedTime.Date()

	fmt.Printf("Last updated at %d-%d-%d\n\n", year, month, day)

	if milestone.Description != "" {
		fmt.Println(milestone.Description)
	} else {
		fmt.Println("No description")
	}
}
