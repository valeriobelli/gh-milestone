package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestones/internal/pkg/application/view"
)

var web bool

var viewMilestoneShort string = `Display the milestone informations. 

With '--web', open the pull request in a web browser instead.
`

var viewCommand = &cobra.Command{
	Use:   "view",
	Short: viewMilestoneShort,
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
			return errors.New("A numeric identifier is needed to view a milestone's information")
		}

		return nil
	},
}

func init() {
	viewCommand.Flags().BoolVarP(&web, "web", "w", false, "View milestone on the browser")

	rootCommand.AddCommand(viewCommand)
}
