package constants

const DateFormat = "2006-01-02"

const (
	MilestoneStateAll    string = "ALL"
	MilestoneStateClosed string = "CLOSED"
	MilestoneStateOpen   string = "OPEN"
)

var MilestoneStates = []string{MilestoneStateAll, MilestoneStateClosed, MilestoneStateOpen}
