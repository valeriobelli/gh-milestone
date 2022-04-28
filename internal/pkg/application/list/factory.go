package list

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/savaki/jq"
	"github.com/shurcooL/githubv4"

	github_entities "github.com/valeriobelli/gh-milestones/internal/pkg/domain/github"
	"github.com/valeriobelli/gh-milestones/internal/pkg/infrastructure/gh"
	"github.com/valeriobelli/gh-milestones/internal/pkg/infrastructure/github"
	"github.com/valeriobelli/gh-milestones/internal/pkg/infrastructure/http"
	tw "github.com/valeriobelli/gh-milestones/internal/pkg/infrastructure/tableWriter"
)

var query struct {
	Repository struct {
		Milestones struct {
			Nodes []github_entities.Milestone
		} `graphql:"milestones(first: $first, orderBy: $orderBy, query: $query, states: $states)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

type MilestonesOrderBy struct {
	Direction string
	Field     string
}

type ListMilestonesConfig struct {
	First   int
	Jq      string
	OrderBy MilestonesOrderBy
	Output  string
	Query   string
	States  []string
}

type ListMilestones struct {
	config ListMilestonesConfig
}

func NewListMilestones(config ListMilestonesConfig) *ListMilestones {
	return &ListMilestones{config: config}
}

func (l ListMilestones) Execute() {
	repoInfo, err := gh.RetrieveRepoInformation()

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	client := github.NewGraphQlClient(http.NewClient())

	err = client.Query(context.Background(), &query, map[string]interface{}{
		"first": githubv4.Int(l.config.First),
		"name":  githubv4.String(strings.TrimSpace(repoInfo.Name)),
		"orderBy": githubv4.MilestoneOrder{
			Direction: githubv4.OrderDirection(l.config.OrderBy.Direction),
			Field:     githubv4.MilestoneOrderField(l.config.OrderBy.Field),
		},
		"owner":  githubv4.String(strings.TrimSpace(repoInfo.Owner)),
		"query":  githubv4.String(l.config.Query),
		"states": l.getStates(),
	})

	if err != nil {
		fmt.Printf(err.Error())

		return
	}

	milestones := query.Repository.Milestones.Nodes

	switch l.config.Output {
	case "json":
		l.printMilestonesAsJson(milestones)
	case "table":
		fallthrough
	default:
		l.printMilestonesAsTable(milestones)
	}
}

func (l ListMilestones) printMilestonesAsTable(milestones []github_entities.Milestone) {
	if len(milestones) == 0 {
		fmt.Println("No milestones found!")

		return
	}

	rows := [][]string{}

	for _, milestone := range milestones {
		rows = append(rows, []string{
			l.printColoredNumber(milestone),
			milestone.Title,
			milestone.Url,
		})
	}

	header := []string{fmt.Sprintf("Showing %d milestones", len(milestones))}
	centerSeparator := ""
	columnSeparator := ""
	rowSeparator := ""
	tablePadding := "\t"

	tw.NewTableWriter(os.Stdout, tw.TableWriterConfig{
		CenterSeparator: &centerSeparator,
		ColumnSeparator: &columnSeparator,
		Header:          &header,
		RowSeparator:    &rowSeparator,
		TablePadding:    &tablePadding,
	}).RenderTable(rows)
}

func (l ListMilestones) printMilestonesAsJson(milestones []github_entities.Milestone) {
	jsonOutput, err := json.Marshal(milestones)

	if err != nil {
		fmt.Print(err.Error())

		return
	}

	jq, err := jq.Parse(l.config.Jq)

	if l.config.Jq == "" || err != nil {
		fmt.Println(string(jsonOutput))

		return
	}

	jqOutput, err := jq.Apply(jsonOutput)

	if err != nil {
		fmt.Print(err.Error())

		return
	}

	fmt.Println(strings.TrimSpace(string(jqOutput)))
}

func (l ListMilestones) printColoredNumber(milestone github_entities.Milestone) string {
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	if milestone.Closed {
		return red.Sprintf("#%d", milestone.Number)
	}

	return green.Sprintf("#%d", milestone.Number)
}

func (l ListMilestones) getStates() []githubv4.MilestoneState {
	mappedSlice := make([]githubv4.MilestoneState, len(l.config.States))

	for idx, state := range l.config.States {
		mappedSlice[idx] = githubv4.MilestoneState(state)
	}

	return mappedSlice
}
