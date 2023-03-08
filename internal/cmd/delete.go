package cmd

import (
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/delete"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/cmdutil"
)

func newDeleteCommand() *cobra.Command {
	deleteCommand := &cobra.Command{
		Use: "delete",
		Example: heredoc.Doc(`
			# autoconfirm the deletion
			$ gh milestone delete 42 --confirm	
		`),
		Long:  "Delete a milestone",
		Short: "Delete a milestone",
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) < 1 {
				return command.Help()
			}

			milestoneId, _ := strconv.Atoi(args[0])

			confirm, _ := command.Flags().GetBool("confirm")
			repo, _ := command.Parent().PersistentFlags().GetString("repo")

			return delete.NewDeleteMilestone(delete.DeleteMilestoneConfig{
				Confirm: confirm,
				Repo:    repo,
			}).Execute(milestoneId)
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

	deleteCommand.SetHelpFunc(cmdutil.HelpFunction)
	deleteCommand.SetUsageFunc(cmdutil.UsageFunction)

	deleteCommand.Flags().BoolP("confirm", "c", false, "Confirm deletion without prompting")

	return deleteCommand
}
