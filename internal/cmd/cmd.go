package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/cmdutil"
)

const version = "v1.1.1"

func Execute() {
	var rootCommand = &cobra.Command{
		Use:           "milestone",
		Short:         "Manage with Github Milestones",
		Long:          "Work with Github milestones.",
		SilenceErrors: true,
		RunE: func(command *cobra.Command, args []string) error {
			versionFlag, _ := command.Flags().GetBool("version")

			if versionFlag {
				fmt.Println(version)

				return nil
			}

			return command.Help()
		},
	}

	rootCommand.SetHelpFunc(cmdutil.HelpFunction)
	rootCommand.SetUsageFunc(cmdutil.UsageFunction)

	rootCommand.Flags().BoolP("version", "v", false, "Print the version of this extension")

	rootCommand.AddCommand(newCreateCommand())
	rootCommand.AddCommand(newDeleteCommand())
	rootCommand.AddCommand(newEditCommand())
	rootCommand.AddCommand(newListCommand())
	rootCommand.AddCommand(newViewCommand())

	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
