package commands_edit

import (
	"fmt"
	"strings"

	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/slices"
)

type stateFlag struct{ *string }

func (flag *stateFlag) String() string {
	if flag.string == nil || *flag.string == "" {
		return ""
	}

	return *flag.string
}

func (flag *stateFlag) GetValue() *string {
	if flag.string == nil || *flag.string == "" {
		return nil
	}

	return flag.string
}

func (flag *stateFlag) Set(value string) error {
	if slices.Contains(constants.CreateMilestoneStates, strings.ToUpper(value)) {
		loweredValue := strings.ToLower(value)

		*flag = stateFlag{string: &loweredValue}

		return nil
	}

	return fmt.Errorf("must be one of {%s}", constants.JoinedEditMilestoneStates)
}

func (flag *stateFlag) Type() string {
	return "state"
}

func NewStateFlag() *stateFlag {
	return &stateFlag{string: nil}
}
