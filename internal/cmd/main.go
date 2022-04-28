package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "v0.1.2"

var rootCommand = &cobra.Command{
	Use:   "milestones",
	Short: "Work with Github Milestons",
	Long:  "A gh extension for viewing and manipulating milestones directly from the terminal.",
	Run: func(command *cobra.Command, args []string) {
		versionFlag, _ := command.Flags().GetBool("version")

		if versionFlag {
			fmt.Println(version)

			return
		}

		command.Help()
	},
}

func init() {
	rootCommand.Flags().BoolP("version", "v", false, "Print the version of this extension")
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
