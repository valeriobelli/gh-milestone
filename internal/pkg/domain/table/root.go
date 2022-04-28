package table

import (
	"fmt"
	"strconv"
)

type TableRow struct {
	Description string
	Number      int
	Title       string
	State       string
}

func (tableRow TableRow) GetDescription() string {
	if tableRow.Description == "" {
		return "<No description>"
	}

	return tableRow.Description
}

func (tableRow TableRow) GetNumber() string {
	return fmt.Sprintf("#%s", strconv.Itoa(tableRow.Number))
}

func (tableRow TableRow) GetState() string {
	return tableRow.State
}

func (tableRow TableRow) GetTitle() string {
	return tableRow.Title
}
