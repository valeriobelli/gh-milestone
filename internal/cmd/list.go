package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/list"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/slices"
)

func NewListCommand() *cobra.Command {
	var getOrderBy = func(command *cobra.Command) list.MilestonesOrderBy {
		orderByField, err := command.Flags().GetString("orderBy.field")
		orderByDirection, err := command.Flags().GetString("orderBy.direction")

		if err != nil {
			return list.MilestonesOrderBy{Direction: "ASC", Field: "NUMBER"}
		}

		return list.MilestonesOrderBy{Direction: orderByDirection, Field: orderByField}
	}

	var getState = func(command *cobra.Command) string {
		state, err := command.Flags().GetString("state")

		if err != nil {
			return constants.MilestoneStateOpen
		}

		uppercaseState := strings.ToUpper(state)

		if slices.Contains(constants.MilestoneStates, uppercaseState) {
			return uppercaseState
		}

		return constants.MilestoneStateOpen
	}

	listCommand := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the available milestones",
		Long: heredoc.Doc(
			`List the available milestones on Github. 
	
			Optionally, the Milestones can be filtered by a search string and status and ordered by some criterias.

			This command permit to print the output as a JSON string and interact with this latter using jq.
		`),
		Run: func(command *cobra.Command, args []string) {
			first, _ := command.Flags().GetInt("first")
			orderBy := getOrderBy(command)
			output, _ := command.Flags().GetString("output")
			query, _ := command.Flags().GetString("query")
			jq, _ := command.Flags().GetString("jq")
			state := getState(command)

			list.NewListMilestones(list.ListMilestonesConfig{
				Query:   query,
				OrderBy: orderBy,
				First:   first,
				Jq:      jq,
				Output:  output,
				State:   state,
			}).Execute()
		},
	}

	possibleStateValues := strings.Join(slices.Map(constants.MilestoneStates, func(value string) string {
		return strings.ToLower(value)
	}), "|")

	listCommand.Flags().IntP("first", "f", 100, "View the first n elements from the list")
	listCommand.Flags().StringP("state", "s", "open", fmt.Sprintf("View milestones by their state: {%s}", possibleStateValues))
	listCommand.Flags().String("orderBy.direction", "ASC", "Use the defined sorting direction")
	listCommand.Flags().String("orderBy.field", "NUMBER", "Sort the milestones by a field ['DUE_DATE', 'CREATED_AT', 'NUMBER', 'UPDATED_AT']")
	listCommand.Flags().StringP("jq", "j", "", "Filter JSON output using a jq expression in combination with --output=json")
	listCommand.Flags().StringP("output", "o", "table", "Decide the output of this command ['json', 'table']")
	listCommand.Flags().StringP("query", "q", "", "View milestones by string pattern")

	return listCommand
}

func init() {
	rootCommand.AddCommand(NewListCommand())
}
