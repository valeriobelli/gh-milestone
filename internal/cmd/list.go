package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/list"
	commands_list "github.com/valeriobelli/gh-milestone/internal/pkg/domain/commands/list"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/cmdutil"
)

func newListCommand() *cobra.Command {
	state := commands_list.NewStateFlag()
	orderByDirection := commands_list.NewOrderByDirectionFlag()
	orderByField := commands_list.NewOrderByFieldFlag()
	output := commands_list.NewOutputFlag()

	listCommand := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the available milestones",
		Example: heredoc.Doc(`
			List closed Milestones
			$ gh milestone list --state closed

			Search by a pattern
			$ gh milestone list --query \"Foo bar\"

			Get first ten milestones
			$ gh milestone list --first 10
		`),
		Long: heredoc.Doc(
			`List the available milestones on Github. 
	
			Optionally, the Milestones can be filtered by a search string and status and ordered by some criterias.

			This command permit to print the output as a JSON string and interact with this latter using jq.
		`),
		RunE: func(command *cobra.Command, args []string) error {
			firstStr, _ := command.Flags().GetString("first")
			first, err := strconv.Atoi(firstStr)

			if err != nil {
				return fmt.Errorf("the value \"%s\" used fort the \"--first\" flag is not a valid number", firstStr)
			}

			query, _ := command.Flags().GetString("query")
			jq, _ := command.Flags().GetString("jq")

			return list.NewListMilestones(list.ListMilestonesConfig{
				Query: query,
				OrderBy: list.MilestonesOrderBy{
					Direction: strings.ToUpper(orderByDirection.String()),
					Field:     strings.ToUpper(orderByField.String()),
				},
				First:  first,
				Jq:     jq,
				Output: output.String(),
				State:  strings.ToUpper(state.String()),
			}).Execute()
		},
	}

	listCommand.SetHelpFunc(cmdutil.HelpFunction)
	listCommand.SetUsageFunc(cmdutil.UsageFunction)

	listCommand.Flags().StringP("first", "f", "100", "View the first N elements from the list")
	listCommand.Flags().VarP(
		state,
		"state",
		"s",
		fmt.Sprintf("View milestones by their state: {%s}", constants.JoinedListMilestoneStates),
	)
	listCommand.Flags().Var(
		orderByDirection,
		"orderBy.direction",
		fmt.Sprintf("Use the defined sorting direction: {%s}", constants.JoinedOrderByDirections),
	)
	listCommand.Flags().Var(
		orderByField,
		"orderBy.field",
		fmt.Sprintf("Sort the milestones by a field: {%s}", constants.JoinedOrderByFields),
	)
	listCommand.Flags().StringP("jq", "j", "", "Filter JSON output using a jq expression in combination with --output=json")
	listCommand.Flags().VarP(
		output,
		"output",
		"o",
		fmt.Sprintf("Decide the output of this command: {%s}", constants.JoinedOutputs),
	)
	listCommand.Flags().StringP("query", "q", "", "View milestones by string pattern")

	return listCommand
}
