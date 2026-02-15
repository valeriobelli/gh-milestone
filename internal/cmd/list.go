package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/valeriobelli/gh-milestone/internal/pkg/application/list"
	commands_list "github.com/valeriobelli/gh-milestone/internal/pkg/domain/commands/list"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/slices"
)

func newListCommand() *cobra.Command {
	state := commands_list.NewStateFlag()
	orderByDirection := commands_list.NewOrderByDirectionFlag()
	orderByField := commands_list.NewOrderByFieldFlag()

	listCommand := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the available milestones",
		Example: heredoc.Doc(`
			List closed Milestones
			$ gh milestone list --state closed

			Search by a pattern
			$ gh milestone list --query "Foo bar"

			Get first ten milestones
			$ gh milestone list --first 10

			Get all milestones
			$ gh milestone list --all

			Print milestones as JSON
			$ gh milestone list --json id
			$ gh milestone list --json id,progressPercentage --json number

			Access Milestone attributes via jq
			$ gh milestone list --json id,progressPercentage --json number --jq ".[0].id"
		`),
		Long: heredoc.Doc(
			`List the available milestones on Github. 
	
			Optionally, the Milestones can be filtered by a search string and status and ordered by some criteria.

			This command permit to print the output as a JSON string and interact with this latter using jq.
		`),
		PreRunE: func(command *cobra.Command, args []string) error {
			all, _ := command.Flags().GetBool("all")

			if !all {
				if value, err := command.Flags().GetInt("first"); err != nil {
					return err
				} else if value < 1 {
					return fmt.Errorf("invalid value for --first: %v", value)
				}
			}

			json := command.Flags().Lookup("json")
			jq := command.Flags().Lookup("jq")

			if json.Changed {
				for _, field := range json.Value.(pflag.SliceValue).GetSlice() {
					if !slices.Contains(github.MilestoneFields, field) {
						command.SilenceUsage = true

						return fmt.Errorf(
							"Unknown JSON field: %q\nAvailable fields:\n  %s",
							field,
							strings.Join(github.MilestoneFields, "\n  "),
						)
					}
				}

				return nil
			} else if jq.Changed {
				return errors.New("cannot use `--jq` without specifying `--json`")
			}

			return nil
		},
		RunE: func(command *cobra.Command, args []string) error {
			all, _ := command.Flags().GetBool("all")
			first, _ := command.Flags().GetInt("first")
			jsonFields, _ := command.Flags().GetStringSlice("json")
			query, _ := command.Flags().GetString("query")
			jq, _ := command.Flags().GetString("jq")
			repo, _ := command.Parent().PersistentFlags().GetString("repo")

			return list.NewListMilestones(list.ListMilestonesConfig{
				All:   all,
				Query: query,
				OrderBy: list.MilestonesOrderBy{
					Direction: strings.ToUpper(orderByDirection.String()),
					Field:     strings.ToUpper(orderByField.String()),
				},
				First: first,
				Jq:    jq,
				Json:  jsonFields,
				Repo:  repo,
				State: strings.ToUpper(state.String()),
			}).Execute()
		},
	}

	listCommand.SetFlagErrorFunc(func(command *cobra.Command, err error) error {
		if command == listCommand && err.Error() == "flag needs an argument: --json" {
			command.SilenceUsage = true

			sort.Strings(github.MilestoneFields)

			return fmt.Errorf(
				"Specify one or more comma-separated fields for `--json`:\n  %s",
				strings.Join(github.MilestoneFields, "\n  "),
			)
		}

		if command.HasParent() {
			return command.Parent().FlagErrorFunc()(command, err)
		}

		return err
	})

	listCommand.Flags().BoolP("all", "a", false, "Retrieve all milestones")
	listCommand.Flags().IntP("first", "f", 100, "View the first N elements from the list")
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
	listCommand.Flags().StringP("jq", "j", "", "Filter JSON output using a jq expression in combination with --json")
	listCommand.Flags().StringSlice("json", nil, "Output JSON with the specified fields")
	listCommand.Flags().StringP("query", "q", "", "View milestones by string pattern")

	return listCommand
}
