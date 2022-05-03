package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/view"
)

func NewViewCommand() *cobra.Command {
	viewCommand := &cobra.Command{
		Use:   "view",
		Short: "Display the milestone",
		Long: heredoc.Doc(`
			Display the information of the milestone.

			Optionally, to open the pull request in a browser the '--web' flag can be passed.
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

			web, _ := command.Flags().GetBool("web")

			view.NewViewMilestone(view.ViewMilestoneConfig{
				Web: web,
			}).Execute(milestoneNumber)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return nil
			}

			_, err := strconv.Atoi(args[0])

			if err != nil {
				return errors.New("A numeric identifier is needed to view the milestone")
			}

			return nil
		},
	}

	viewCommand.Flags().BoolP("web", "w", false, "View milestone on the browser")

	return viewCommand
}

func init() {
	rootCommand.AddCommand(NewViewCommand())
}
