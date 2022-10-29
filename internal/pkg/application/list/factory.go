package list

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/shurcooL/githubv4"

	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	github_entities "github.com/valeriobelli/gh-milestone/internal/pkg/domain/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/gh"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/http"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/spinner"
	tw "github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/tableWriter"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/jq"
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
	Json    []string
	Query   string
	State   string
}

type ListMilestones struct {
	config ListMilestonesConfig
}

func NewListMilestones(config ListMilestonesConfig) *ListMilestones {
	return &ListMilestones{config: config}
}

func (l ListMilestones) Execute() error {
	repoInfo, err := gh.RetrieveRepoInformation()

	if err != nil {
		return err
	}

	client := github.NewGraphQlClient(http.NewClient())

	spinner := spinner.NewSpinner()

	spinner.Start()

	err = client.Query(context.Background(), &query, map[string]interface{}{
		"first": githubv4.Int(l.config.First),
		"name":  githubv4.String(strings.TrimSpace(repoInfo.Name)),
		"orderBy": githubv4.MilestoneOrder{
			Direction: githubv4.OrderDirection(strings.ToUpper(l.config.OrderBy.Direction)),
			Field:     githubv4.MilestoneOrderField(strings.ToUpper(l.config.OrderBy.Field)),
		},
		"owner":  githubv4.String(strings.TrimSpace(repoInfo.Owner)),
		"query":  githubv4.String(l.config.Query),
		"states": l.getStates(),
	})

	spinner.Stop()

	if err != nil {
		return err
	}

	milestones := query.Repository.Milestones.Nodes

	if len(l.config.Json) > 0 {
		return l.printMilestonesAsJson(l.config.Json, milestones)
	}

	return l.printMilestonesAsTable(milestones)
}

func (l ListMilestones) printMilestonesAsTable(milestones []github_entities.Milestone) error {
	if len(milestones) == 0 {
		fmt.Println("No milestones found!")

		return nil
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

	return nil
}

func (l ListMilestones) printMilestonesAsJson(jsonFields []string, milestones []github_entities.Milestone) error {
	data, err := json.Marshal(milestones)

	if err != nil {
		return err
	}

	var unmarshaledMilestones []map[string]interface{}

	if err = json.Unmarshal(data, &unmarshaledMilestones); err != nil {
		return err
	}

	milestoneWithFilteredProperties := []interface{}{}

	for _, milestone := range unmarshaledMilestones {
		mappedMilestone := map[string]interface{}{}

		for _, field := range jsonFields {
			mappedMilestone[field] = milestone[field]
		}

		milestoneWithFilteredProperties = append(milestoneWithFilteredProperties, mappedMilestone)
	}

	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(milestoneWithFilteredProperties); err != nil {
		return err
	}

	return jq.Evaluate(&buf, os.Stdout, l.config.Jq)
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
	if l.config.State == constants.MilestoneStateAll {
		return []githubv4.MilestoneState{githubv4.MilestoneStateOpen, githubv4.MilestoneStateClosed}
	}

	return []githubv4.MilestoneState{githubv4.MilestoneState(l.config.State)}
}
