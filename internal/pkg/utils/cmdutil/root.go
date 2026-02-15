package cmdutil

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/slices"
)

func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)

	return fmt.Sprintf(template, s)
}

func HelpFunction(command *cobra.Command, args []string) {
	var getCommandDescription = func() string {
		if command.Long != "" {
			return command.Long
		}

		return command.Short
	}

	bold := color.New(color.Bold)

	fmt.Println(getCommandDescription())
	fmt.Println()

	_, _ = bold.Println("USAGE")
	fmt.Printf("  gh %s\n\n", command.UseLine())

	commands := command.Commands()

	if len(commands) > 0 {
		_, _ = bold.Println("CORE COMMANDS")

		for _, command := range commands {
			fmt.Printf("  %s%s\n", rpad(command.Name()+":", 12), command.Short)
		}

		fmt.Println()
	}

	_, _ = bold.Println("FLAGS")
	fmt.Print(command.LocalFlags().FlagUsages())

	fmt.Println()

	_, _ = bold.Println("INHERITED FLAGS")
	fmt.Print(command.InheritedFlags().FlagUsages())

	if command.Example != "" {
		fmt.Println()
		_, _ = bold.Println("EXAMPLES")
		fmt.Println(
			strings.Join(
				slices.Map(
					strings.Split(command.Example, "\n"),
					func(line string) string {
						return fmt.Sprintf("  %s", line)
					},
				),
				"\n",
			),
		)
	}
}

func UsageFunction(command *cobra.Command) error {
	fmt.Printf("Usage: gh %s\n\n", command.UseLine())

	fmt.Printf("Flags:\n")
	fmt.Print(command.LocalFlags().FlagUsages())

	return nil
}
