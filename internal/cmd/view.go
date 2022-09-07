package cmd

import (
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/view"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/cmdutil"
)

func newViewCommand() *cobra.Command {
	viewCommand := &cobra.Command{
		Use: "view",
		Example: heredoc.Doc(`
			# view the milestone on the browser
			$ gh milestone view 42 --web
		`),
		Long:  "Display a milestone",
		Short: "Display a milestone",
		RunE: func(command *cobra.Command, args []string) error {
			if len(args) < 1 {
				return command.Help()
			}

			milestoneNumber, _ := strconv.Atoi(args[0])

			web, _ := command.Flags().GetBool("web")

			return view.NewViewMilestone(view.ViewMilestoneConfig{
				Web: web,
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

	viewCommand.SetHelpFunc(cmdutil.HelpFunction)
	viewCommand.SetUsageFunc(cmdutil.UsageFunction)

	viewCommand.Flags().BoolP("web", "w", false, "View milestone on the browser")

	return viewCommand
}
