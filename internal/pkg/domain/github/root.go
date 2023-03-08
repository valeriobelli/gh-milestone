package github

import (
	"fmt"
	"time"

	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
)

type Milestone struct {
	Closed             bool    `json:"closed"`
	Description        string  `json:"description"`
	DueOn              string  `json:"dueOn"`
	Id                 string  `json:"id"`
	Number             int     `json:"number"`
	ProgressPercentage float64 `json:"progressPercentage"`
	State              string  `json:"state"`
	Title              string  `json:"title"`
	UpdatedAt          string  `json:"updatedAt"`
	Url                string  `json:"url"`
}

var MilestoneFields = []string{
	"closed",
	"description",
	"dueOn",
	"id",
	"number",
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
