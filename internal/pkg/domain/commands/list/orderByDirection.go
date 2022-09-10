package commands_list

import (
	"fmt"
	"strings"

	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/slices"
)

type orderByDirectionFlag struct{ *string }

func (flag *orderByDirectionFlag) String() string {
	if flag.string == nil || *flag.string == "" {
		return strings.ToLower(constants.OrderByDirectionAsc)
	}

	return *flag.string
}

func (flag *orderByDirectionFlag) Set(value string) error {
	if slices.Contains(constants.OrderByDirections, strings.ToUpper(value)) {
		*flag = orderByDirectionFlag{string: &value}

		return nil
	}

	return fmt.Errorf("must be one of {%s}", constants.JoinedOrderByDirections)
}

func (flag *orderByDirectionFlag) Type() string {
	return "orderBy.direction"
}

func NewOrderByDirectionFlag() *orderByDirectionFlag {
	return &orderByDirectionFlag{string: nil}
}
