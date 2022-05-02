package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/edit"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
)

func NewEditCommand() *cobra.Command {
	var getDescription = func(command *cobra.Command) *string {
		description, _ := command.Flags().GetString("description")

		if description == "" {
			return nil
		}

		return &description
	}

	var parseDueOn = func(dueOn string) (*time.Time, error) {
		parsedDueDate, err := time.Parse(constants.DateFormat, dueOn)

		if err != nil {
			return nil, err
		}

		return &parsedDueDate, nil
	}

	var getDueDate = func(command *cobra.Command) (*time.Time, error) {
		dueDate, err := command.Flags().GetString("due-date")

		if err != nil {
			return nil, err
		}

		if dueDate == "" {
			return nil, nil
		}

		return parseDueOn(dueDate)
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
		Short: "Edit Github Milestones",
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
			verbose, _ := command.Flags().GetBool("verbose")

			if err != nil {
				fmt.Println(err.Error())

				return
			}

			if description == nil && dueDate == nil && state == nil && title == nil {
				fmt.Println("At least on option among ['description', 'dueOn', 'state', 'title'] is needed to edit the milestone")

			}

			edit.NewEditMilestone(edit.EditMilestoneConfig{
				Description: description,
				DueDate:     dueDate,
				State:       state,
				Title:       title,
				Verbose:     verbose,
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

	editCommand.Flags().String("due-date", "", fmt.Sprintf("Set the Due date is %s", constants.DateFormat))
	editCommand.Flags().String("description", "", "Milestone's description")
	editCommand.Flags().String("state", "", "Milestone's state. Accepted values are: ['open', 'closed']")
	editCommand.Flags().String("title", "", "Milestone's title")
	editCommand.Flags().BoolP("verbose", "v", false, "Print the result of the editing")

	return editCommand
}

func init() {
	rootCommand.AddCommand(NewEditCommand())
}
