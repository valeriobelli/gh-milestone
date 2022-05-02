package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/edit"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/github"
)

func NewEditCommand() *cobra.Command {
	var getDescription = func(command *cobra.Command) *string {
		description, _ := command.Flags().GetString("description")

		if description == "" {
			return nil
		}

		return &description
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

	var getState = func(command *cobra.Command) *string {
		state, _ := command.Flags().GetString("state")

		if state == "" || (state != "open" && state != "closed") {
			return nil
		}

		return &state
	}

	var getTitle = func(command *cobra.Command) *string {
		title, _ := command.Flags().GetString("title")

		if title == "" {
			return nil
		}

		return &title
	}

	editCommand := &cobra.Command{
		Use:   "edit",
		Short: "Edit a milestone",
		Long:  "Edit a milestone on Github.",
		Example: heredoc.Doc(`
			# set a new title
			gh milestone edit 1 --title "<new title>"

			# close a milestone
			gh milestone edit 1 --state closed
		`),
		Run: func(command *cobra.Command, args []string) {
			if len(args) == 0 {
				command.Help()

				return
			}

			milestoneNumber, err := strconv.Atoi(args[0])

			if err != nil {
				fmt.Println(err.Error())

				return
			}

			description := getDescription(command)
			dueDate, err := getDueDate(command)
			state := getState(command)
			title := getTitle(command)

			if err != nil {
				fmt.Println(err.Error())

				return
			}

			edit.NewEditMilestone(edit.EditMilestoneConfig{
				Description: description,
				DueDate:     dueDate,
				State:       state,
				Title:       title,
			}).Execute(milestoneNumber)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return nil
			}

			_, err := strconv.Atoi(args[0])

			if err != nil {
				return errors.New("An numeric identifier is needed to edit a Milestone")
			}

			return nil
		},
	}

	editCommand.Flags().StringP("due-date", "u", "", fmt.Sprintf("Set the milestone Due date [%s]", constants.DateFormat))
	editCommand.Flags().StringP("description", "d", "", "Set the milestone description")
	editCommand.Flags().StringP("state", "s", "", "Set the milestone state ['open', 'closed']")
	editCommand.Flags().StringP("title", "t", "", "Set the milestone title")

	return editCommand
}

func init() {
	rootCommand.AddCommand(NewEditCommand())
}
