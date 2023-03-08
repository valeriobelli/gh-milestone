package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/create"
	commands_create "github.com/valeriobelli/gh-milestone/internal/pkg/domain/commands/create"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/cmdutil"
)

func newCreateCommand() *cobra.Command {
	dueDate := commands_create.NewDueDateFlag()

	createCommand := &cobra.Command{
		Use:   "create",
		Short: "Create a milestone",
		Example: heredoc.Doc(`
			# create a new milestone with a title by using flags
			$ gh milestone create --title v1.0.0

			# create a new milestone by using flags
			gh milestone create --title v1.0.0 --description "# This is a description" --due-date 2022-06-01

			# create a new milestone interactively
			gh milestone create
		`),
		Long: heredoc.Doc(`
			Create a milestone on Github.

			Optionally, this command permits to create a Milestone interactively when flags 
			for required fields are not defined.
		`),
		RunE: func(command *cobra.Command, args []string) error {
			dueDate, err := dueDate.GetValue()

			if err != nil {
				return err
			}

			description, _ := command.Flags().GetString("description")
			title, _ := command.Flags().GetString("title")
			repo, _ := command.Parent().PersistentFlags().GetString("repo")

			return create.NewCreateMilestone(create.CreateMilestoneConfig{
				Description: description,
				DueDate:     dueDate,
				Repo:        repo,
				Title:       title,
			}).Execute()
		},
	}

	createCommand.SetHelpFunc(cmdutil.HelpFunction)
	createCommand.SetUsageFunc(cmdutil.UsageFunction)

	createCommand.Flags().StringP("description", "d", "", "Set the description")
	createCommand.Flags().VarP(dueDate, "due-date", "u", fmt.Sprintf("Set the due date [%s]", constants.DateFormat))
	createCommand.Flags().StringP("title", "t", "", "Set the title")

	return createCommand
}
