package github

import (
	"fmt"
	"time"

	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
)

type MilestoneIssueCount struct {
	TotalCount int `json:"totalCount"`
}

type Milestone struct {
	Closed             bool                `json:"closed"`
	ClosedIssues       MilestoneIssueCount `json:"closedIssues" graphql:"closedIssues: issues(states: CLOSED)"`
	Description        string              `json:"description"`
	DueOn              string              `json:"dueOn"`
	Id                 string              `json:"id"`
	Issues             MilestoneIssueCount `json:"issues"`
	Number             int                 `json:"number"`
	OpenIssues         MilestoneIssueCount `json:"openIssues" graphql:"openIssues: issues(states: OPEN)"`
	ProgressPercentage float64             `json:"progressPercentage"`
	State              string              `json:"state"`
	Title              string              `json:"title"`
	UpdatedAt          string              `json:"updatedAt"`
	Url                string              `json:"url"`
}

var MilestoneFields = []string{
	"closed",
	"closedIssues",
	"description",
	"dueOn",
	"id",
	"issues",
	"number",
	"openIssues",
	"progressPercentage",
	"state",
	"title",
	"updatedAt",
	"url",
}

type DueDate struct{ time.Time }

func NewDueDate(dueDate string) (*DueDate, error) {
	currentTime := time.Now()

	parsedDueDate, err := time.ParseInLocation(constants.DateFormat, dueDate, currentTime.Location())

	if err != nil {
		return nil, fmt.Errorf("due date %s is not respecting the format %s", dueDate, constants.DateFormat)
	}

	return &DueDate{parsedDueDate}, nil
}
