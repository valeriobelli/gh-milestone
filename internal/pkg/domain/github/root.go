package github

import (
	"errors"
	"fmt"
	"time"

	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
)

type Milestone struct {
	Closed             bool
	Description        string
	DueOn              string
	Id                 string
	Number             int
	ProgressPercentage float64
	State              string
	Title              string
	Url                string
	UpdatedAt          string
}

type DueDate struct{ time.Time }

func NewDueDate(dueDate string) (*DueDate, error) {
	parsedDueDate, err := time.Parse(constants.DateFormat, dueDate)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Due date %s is not respecting the format %s", dueDate, constants.DateFormat))
	}

	return &DueDate{parsedDueDate}, nil
}
