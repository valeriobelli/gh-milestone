package commands_list

import (
	"fmt"
	"strings"

	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/constants"
	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/slices"
)

type orderByFieldFlag struct{ *string }

func (flag *orderByFieldFlag) String() string {
	if flag.string == nil || *flag.string == "" {
		return strings.ToLower(constants.OrderByFieldNumber)
	}

	return *flag.string
}

func (flag *orderByFieldFlag) Set(value string) error {
	if slices.Contains(constants.OrderByFields, strings.ToUpper(value)) {
		*flag = orderByFieldFlag{string: &value}

		return nil
	}

	return fmt.Errorf("must be one of {%s}", constants.JoinedOrderByFields)
}

func (flag *orderByFieldFlag) Type() string {
	return "orderBy.field"
}

func NewOrderByFieldFlag() *orderByFieldFlag {
	return &orderByFieldFlag{string: nil}
}
