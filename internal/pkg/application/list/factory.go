package list

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
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
			Nodes    []github_entities.Milestone
			PageInfo struct {
				HasNextPage bool
				EndCursor   githubv4.String
			}
		} `graphql:"milestones(first: $first, after: $after, orderBy: $orderBy, query: $query, states: $states)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

type MilestonesOrderBy struct {
	Direction string
	Field     string
}

type ListMilestonesConfig struct {
	All     bool
	First   int
	Jq      string
	OrderBy MilestonesOrderBy
	Json    []string
	Query   string
	Repo    string
	State   string
}

type ListMilestones struct {
	config ListMilestonesConfig
}

func NewListMilestones(config ListMilestonesConfig) *ListMilestones {
	return &ListMilestones{config: config}
}

const maxPageSize = 100

func (l ListMilestones) Execute() error {
	repoInfo, err := gh.RetrieveRepoInformation(l.config.Repo)

	if err != nil {
		return err
	}

	client := github.NewGraphQlClient(http.NewClient())

	spinner := spinner.NewSpinner()

	spinner.Start()

	apiOrderBy := l.getApiOrderBy()

	var allMilestones []github_entities.Milestone
	var cursor *githubv4.String
	remaining := l.config.First

	for {
		pageSize := maxPageSize
		if !l.config.All && remaining < pageSize {
			pageSize = remaining
		}

		variables := map[string]interface{}{
			"first": githubv4.Int(pageSize),
			"after": cursor,
			"name":  githubv4.String(strings.TrimSpace(repoInfo.Name)),
			"orderBy": githubv4.MilestoneOrder{
				Direction: githubv4.OrderDirection(strings.ToUpper(apiOrderBy.Direction)),
				Field:     githubv4.MilestoneOrderField(strings.ToUpper(apiOrderBy.Field)),
			},
			"owner":  githubv4.String(strings.TrimSpace(repoInfo.Owner)),
			"query":  githubv4.String(l.config.Query),
			"states": l.getStates(),
		}

		err = client.Query(context.Background(), &query, variables)

		if err != nil {
			spinner.Stop()
			return err
		}

		allMilestones = append(allMilestones, query.Repository.Milestones.Nodes...)

		if !query.Repository.Milestones.PageInfo.HasNextPage {
			break
		}

		if !l.config.All {
			remaining -= len(query.Repository.Milestones.Nodes)
			if remaining <= 0 {
				break
			}
		}

		endCursor := query.Repository.Milestones.PageInfo.EndCursor
		cursor = &endCursor
	}

	spinner.Stop()

	milestones := allMilestones

	if l.needsClientSideSort() {
		milestones = l.sortMilestones(milestones)
	}

	if len(l.config.Json) > 0 {
		return l.printMilestonesAsJson(l.config.Json, milestones)
	}

	return l.printMilestonesAsTable(milestones)
}

func (l ListMilestones) needsClientSideSort() bool {
	field := strings.ToUpper(l.config.OrderBy.Field)

	return field == constants.OrderByFieldTitle || field == constants.OrderByFieldIssues
}

func (l ListMilestones) getApiOrderBy() MilestonesOrderBy {
	if l.needsClientSideSort() {
		return MilestonesOrderBy{
			Direction: constants.OrderByDirectionAsc,
			Field:     constants.OrderByFieldNumber,
		}
	}

	return l.config.OrderBy
}

func (l ListMilestones) sortMilestones(milestones []github_entities.Milestone) []github_entities.Milestone {
	sorted := make([]github_entities.Milestone, len(milestones))
	copy(sorted, milestones)

	ascending := strings.ToUpper(l.config.OrderBy.Direction) != constants.OrderByDirectionDesc
	field := strings.ToUpper(l.config.OrderBy.Field)

	sort.Slice(sorted, func(i, j int) bool {
		switch field {
		case constants.OrderByFieldIssues:
			if ascending {
				return sorted[i].Issues.TotalCount < sorted[j].Issues.TotalCount
			}

			return sorted[i].Issues.TotalCount > sorted[j].Issues.TotalCount
		default:
			if ascending {
				return strings.ToLower(sorted[i].Title) < strings.ToLower(sorted[j].Title)
			}

			return strings.ToLower(sorted[i].Title) > strings.ToLower(sorted[j].Title)
		}
	})

	return sorted
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
