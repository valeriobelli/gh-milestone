package commands_create

import (
	"fmt"
	"time"

	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/github"
)

type dueDateFlag struct{ *string }

func (flag *dueDateFlag) String() string {
	if flag.string == nil || *flag.string == "" {
		return ""
	}

	return *flag.string
}

func (flag *dueDateFlag) GetValue() (*time.Time, error) {
	if flag.string == nil || *flag.string == "" {
		return nil, nil
	}

	parsedDate, err := github.NewDueDate(*flag.string)

	if err != nil {
		return nil, fmt.Errorf(
			"the value \"%s\" is not a valid date that respect the format \"%s\"",
			*flag.string,
			constants.DateFormat,
		)
	}

	return &parsedDate.Time, nil
}

func (flag *dueDateFlag) Set(value string) error {
	*flag = dueDateFlag{string: &value}

	return nil
}

func (flag *dueDateFlag) Type() string {
	return "dueDate"
}

func NewDueDateFlag() *dueDateFlag {
	return &dueDateFlag{string: nil}
}
