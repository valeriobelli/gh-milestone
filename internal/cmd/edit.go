package cmd

import (
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/edit"
	commands_edit "github.com/valeriobelli/gh-milestone/internal/pkg/domain/commands/edit"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/cmdutil"
)

func newEditCommand() *cobra.Command {
	description := commands_edit.NewDescriptionFlag()
	dueDate := commands_edit.NewDueDateFlag()
	state := commands_edit.NewStateFlag()
	title := commands_edit.NewTitleFlag()

	editCommand := &cobra.Command{
		Use: "edit",
		Example: heredoc.Doc(`
			# set the new title
			$ gh milestone edit 42 --title "<new title>"

			# close the milestone
			$ gh milestone edit 42 --state closed
		`),
		Long:  "Edit a milestone",
		Short: "Edit a milestone",
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) < 1 {
				return command.Help()
			}

			milestoneNumber, _ := strconv.Atoi(args[0])

			dueDate, err := dueDate.GetValue()

			if err != nil {
				return err
			}

			description := description.GetValue()
			state := state.GetValue()
			title := title.GetValue()

			if dueDate == nil && description == nil && state == nil && title == nil {
				fmt.Println("Nothing to edit, exiting.")

				return nil
			}

			return edit.NewEditMilestone(edit.EditMilestoneConfig{
				Description: description,
				DueDate:     dueDate,
				State:       state,
				Title:       title,
			}).Execute(milestoneNumber)
		},
		Args: func(command *cobra.Command, args []string) error {
			if len(args) < 1 {
				return nil
			}

			milestoneId := args[0]

			_, err := strconv.Atoi(milestoneId)

			if err != nil {
				return fmt.Errorf(
					"the value \"%s\" is not a valid numeric identifier needed to edit a Milestone",
					milestoneId,
				)
			}

			return nil
		},
	}

	editCommand.SetHelpFunc(cmdutil.HelpFunction)
	editCommand.SetUsageFunc(cmdutil.UsageFunction)

	editCommand.Flags().VarP(
		dueDate,
		"due-date",
		"u",
		fmt.Sprintf("Set the milestone Due date [%s]", constants.DateFormat),
	)
	editCommand.Flags().VarP(description, "description", "d", "Set the milestone description")
	editCommand.Flags().VarP(
		state,
		"state",
		"s",
		fmt.Sprintf("Set the milestone state: {%s}", constants.JoinedEditMilestoneStates),
	)
	editCommand.Flags().VarP(title, "title", "t", "Set the milestone title")

	return editCommand
}
