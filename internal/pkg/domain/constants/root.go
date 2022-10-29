package constants

import (
	"strings"

	"github.com/valeriobelli/gh-milestone/internal/pkg/utils/slices"
)

func prepareForDocUsage(values []string) string {
	return strings.Join(slices.Map(values, strings.ToLower), "|")
}

const DateFormat = "2006-01-02"

const (
	MilestoneStateAll    string = "ALL"
	MilestoneStateClosed string = "CLOSED"
	MilestoneStateOpen   string = "OPEN"
)

var ListMilestoneStates = []string{MilestoneStateAll, MilestoneStateClosed, MilestoneStateOpen}
var CreateMilestoneStates = []string{MilestoneStateClosed, MilestoneStateOpen}

var JoinedListMilestoneStates = prepareForDocUsage(ListMilestoneStates)
var JoinedEditMilestoneStates = prepareForDocUsage(CreateMilestoneStates)

const (
	OrderByDirectionAsc  string = "ASC"
	OrderByDirectionDesc string = "DESC"
)

var OrderByDirections = []string{OrderByDirectionAsc, OrderByDirectionDesc}

var JoinedOrderByDirections = prepareForDocUsage(OrderByDirections)

const (
	OrderByFieldCreatedAt string = "CREATED_AT"
	OrderByFieldDueDate   string = "DUE_DATE"
	OrderByFieldNumber    string = "NUMBER"
	OrderByFieldUpdatedAt string = "UPDATED_AT"
)

var OrderByFields = []string{OrderByFieldCreatedAt, OrderByFieldDueDate, OrderByFieldNumber, OrderByFieldUpdatedAt}

var JoinedOrderByFields = prepareForDocUsage(OrderByFields)
