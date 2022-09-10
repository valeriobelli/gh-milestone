package commands_list

import (
	"fmt"
	"strings"

	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/slices"
)

type outputFlag struct{ *string }

func (flag *outputFlag) String() string {
	if flag.string == nil || *flag.string == "" {
		return strings.ToLower(constants.OutputTable)
	}

	return *flag.string
}

func (flag *outputFlag) Set(value string) error {
	if slices.Contains(constants.Outputs, value) {
		*flag = outputFlag{string: &value}

		return nil
	}

	return fmt.Errorf("must be one of {%s}", constants.JoinedOutputs)
}

func (flag *outputFlag) Type() string {
	return "output"
}

func NewOutputFlag() *outputFlag {
	return &outputFlag{string: nil}
}
