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
	Url                string  `json:"url"`
	UpdatedAt          string  `json:"updatedAt"`
}

type DueDate struct{ time.Time }

func NewDueDate(dueDate string) (*DueDate, error) {
	parsedDueDate, err := time.Parse(constants.DateFormat, dueDate)

	if err != nil {
		return nil, fmt.Errorf("due date %s is not respecting the format %s", dueDate, constants.DateFormat)
	}

	return &DueDate{parsedDueDate}, nil
}
