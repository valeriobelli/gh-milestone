package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestones/internal/pkg/application/edit"
	"github.com/valeriobelli/gh-milestones/internal/pkg/domain/constants"
)

func getDescription(command *cobra.Command) *string {
	description, _ := command.Flags().GetString("description")

	if description == "" {
		return nil
	}

	return &description
}

func getDueOn(command *cobra.Command) (*time.Time, error) {
	dueOn, err := command.Flags().GetString("dueOn")

	if err != nil {
		return nil, err
	}

	if dueOn == "" {
		return nil, nil
	}

	parsedDueDate, err := time.Parse(constants.DateFormat, dueOn)

	if err != nil {
		return nil, err
	}

	return &parsedDueDate, nil
}

func getState(command *cobra.Command) *string {
	state, _ := command.Flags().GetString("state")

	if state == "" || (state != "open" && state != "closed") {
		return nil
	}

	return &state
}

func getTitle(command *cobra.Command) *string {
	title, _ := command.Flags().GetString("title")

	if title == "" {
		return nil
	}

	return &title
}

var editCommand = &cobra.Command{
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
		dueOn, err := getDueOn(command)
		state := getState(command)
		title := getTitle(command)
		verbose, _ := command.Flags().GetBool("verbose")

		if err != nil {
			fmt.Println(err.Error())

			return
		}

		if description == nil && dueOn == nil && state == nil && title == nil {
			fmt.Println("At least on option among ['description', 'dueOn', 'state', 'title'] is needed to edit the milestone")

			return
		}

		edit.NewEditMilestone(edit.EditMilestoneConfig{
			Description: description,
			DueOn:       dueOn,
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

func init() {
	editCommand.Flags().String("dueOn", "", fmt.Sprintf("Milestone's due date. Accepted format is %s", constants.DateFormat))
	editCommand.Flags().String("description", "", "Milestone's description")
	editCommand.Flags().String("state", "", "Milestone's state. Accepted values are: ['open', 'closed']")
	editCommand.Flags().String("title", "", "Milestone's title")
	editCommand.Flags().BoolP("verbose", "v", false, "Print the result of the editing")

	rootCommand.AddCommand(editCommand)
}
