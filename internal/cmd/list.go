package cmd

import (
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestones/internal/pkg/application/list"
)

var first int
var jq string
var orderBy string
var orderByDirection string
var orderByField string
var output string
var query string
var states []string

func getOrderBy(command *cobra.Command) list.MilestonesOrderBy {
	orderByField, err := command.Flags().GetString("orderBy.field")
	orderByDirection, err := command.Flags().GetString("orderBy.direction")

	if err != nil {
		return list.MilestonesOrderBy{Direction: "ASC", Field: "NUMBER"}
	}

	return list.MilestonesOrderBy{Direction: orderByDirection, Field: orderByField}
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List the available milestones.",
	Long: `List the available milestones. 

Optionally, the Milestones can be filtered by a search string and status and ordered by some criterias.

This command permit to print the output as a JSON string and interact with this latter using jq.
`,
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

func init() {
	listCommand.Flags().IntVarP(&first, "first", "f", 100, "View the first n elements from the list")
	listCommand.Flags().StringArrayVarP(&states, "states", "s", []string{"OPEN"}, "View milestons by states")
	listCommand.Flags().StringVar(&orderByDirection, "orderBy.direction", "ASC", "Milestone's sorting direction")
	listCommand.Flags().StringVar(&orderByField, "orderBy.field", "NUMBER", "Sort the milestones by a field. Available values are ['DUE_DATE', 'CREATED_AT', 'NUMBER', 'UPDATED_AT']")
	listCommand.Flags().StringVarP(&jq, "jq", "j", "", "Filter JSON output using a jq expression. It works in combination with --output=json")
	listCommand.Flags().StringVarP(&output, "output", "o", "table", "Decide the output of this command. Available formats are ['json', 'table']")
	listCommand.Flags().StringVarP(&query, "query", "q", "", "View milestones by string pattern")

	rootCommand.AddCommand(listCommand)
}
