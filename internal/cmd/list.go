package cmd

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/list"
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

	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List the available milestones",
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
			states, _ := command.Flags().GetStringArray("states")

			list.NewListMilestones(list.ListMilestonesConfig{
				Query:   query,
				OrderBy: orderBy,
				First:   first,
				Jq:      jq,
				Output:  output,
				States:  states,
			}).Execute()
		},
	}

	listCommand.Flags().IntP("first", "f", 100, "View the first n elements from the list")
	listCommand.Flags().StringArrayP("states", "s", []string{"OPEN"}, "View milestons by states")
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
