package commands_list

import (
	"fmt"
	"strings"

	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/slices"
)

type stateFlag struct{ *string }

func (flag *stateFlag) String() string {
	if flag.string == nil || *flag.string == "" {
		return strings.ToLower(constants.MilestoneStateOpen)
	}

	return *flag.string
}

func (flag *stateFlag) Set(value string) error {
	if slices.Contains(constants.ListMilestoneStates, strings.ToUpper(value)) {
		*flag = stateFlag{string: &value}

		return nil
	}

	return fmt.Errorf("must be one of {%s}", constants.JoinedListMilestoneStates)
}

func (flag *stateFlag) Type() string {
	return "state"
}

func NewStateFlag() *stateFlag {
	return &stateFlag{string: nil}
}
