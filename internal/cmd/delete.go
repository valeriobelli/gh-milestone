package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/delete"
)

func NewDeleteCommand() *cobra.Command {
	deleteCommand := &cobra.Command{
		Use:   "delete",
		Short: "Delete milestone",
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

			confirm, _ := command.Flags().GetBool("confirm")

			delete.NewDeleteMilestone(delete.DeleteMilestoneConfig{
				Confirm: confirm,
			}).Execute(milestoneNumber)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return nil
			}

			_, err := strconv.Atoi(args[0])

			if err != nil {
				return errors.New("A numeric identifier is needed to delete the milestone")
			}

			return nil
		},
	}

	deleteCommand.Flags().BoolP("confirm", "c", false, "Confirm deletion without prompting")

	return deleteCommand
}

func init() {
	rootCommand.AddCommand(NewDeleteCommand())
}
