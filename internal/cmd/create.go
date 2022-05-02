package cmd

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/create"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/github"
)

func NewCreateCommand() *cobra.Command {
	var getDescription = func(command *cobra.Command) string {
		description, err := command.Flags().GetString("description")

		if err != nil {
			return ""
		}

		return description
	}

	var getDueDate = func(command *cobra.Command) (*time.Time, error) {
		dueDate, err := command.Flags().GetString("due-date")

		if err != nil {
			return nil, err
		}

		if dueDate == "" {
			return nil, nil
		}

		parsedDate, err := github.NewDueDate(dueDate)

		if err != nil {
			return nil, err
		}

		return &parsedDate.Time, nil
	}

	var getTitle = func(command *cobra.Command) string {
		title, err := command.Flags().GetString("title")

		if err != nil {
			return ""
		}

		return title
	}

	createCommand := &cobra.Command{
		Use:   "create",
		Short: "Create a milestone",
		Long: heredoc.Doc(`
			Create a mileston on Github.

			Optionally, this command permits to create a Milestone interactively when flags 
			for required fields are not defined.
		`),
		Example: heredoc.Doc(`
			# create a new milestone with a title by using flags
			$ gh milestones create --title v1.0.0

			# create a new milestone by using flags
			gh milestones create --title v1.0.0 --description "# This is a description" --due-date 2022-06-01

			# create a new milestone interactively
			gh milestones create
		`),
		Run: func(command *cobra.Command, args []string) {
			description := getDescription(command)
			dueDate, err := getDueDate(command)
			title := getTitle(command)

			if err != nil {
				fmt.Println(err.Error())

				return
			}

			create.NewCreateMilestone(create.CreateMilestoneConfig{
				Description: description,
				DueDate:     dueDate,
				Title:       title,
			}).Execute()
		},
	}

	createCommand.Flags().StringP("description", "d", "", "Set the description")
	createCommand.Flags().StringP("due-date", "u", "", fmt.Sprintf("Set the due date [%s]", constants.DateFormat))
	createCommand.Flags().StringP("title", "t", "", "Set the title")

	return createCommand
}

func init() {
	rootCommand.AddCommand(NewCreateCommand())
}
