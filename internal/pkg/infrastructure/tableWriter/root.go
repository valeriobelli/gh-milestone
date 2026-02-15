package tablewriter

import (
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

type TableWriterConfig struct {
	Border          *bool
	CenterSeparator *string
	ColumnAlignment *[]int
	ColumnSeparator *string
	Header          *[]string
	HeaderLine      *bool
	RowSeparator    *string
	TablePadding    *string
}

func (c TableWriterConfig) getColumnAlignment() []int {
	if c.ColumnAlignment == nil {
		return []int{tablewriter.ALIGN_DEFAULT}
	}

	return *c.ColumnAlignment
}

func (c TableWriterConfig) getCenterSeparator() string {
	if c.CenterSeparator == nil {
		return tablewriter.CENTER
	}

	return *c.CenterSeparator
}

func (c TableWriterConfig) getColumnSeparator() string {
	if c.ColumnSeparator == nil {
		return tablewriter.COLUMN
	}

	return *c.ColumnSeparator
}

func (c TableWriterConfig) getHeader() []string {
	if c.Header == nil {
		return []string{""}
	}

	return *c.Header
}

func (c TableWriterConfig) getRowSeparator() string {
	if c.RowSeparator == nil {
		return tablewriter.ROW
	}

	return *c.RowSeparator
}

func (c TableWriterConfig) getHeaderLine() bool {
	if c.HeaderLine == nil {
		return true
	}

	return *c.HeaderLine
}

func (c TableWriterConfig) getBorder() bool {
	if c.Border == nil {
		return true
	}

	return *c.Border
}

func (c TableWriterConfig) getTablePadding() string {
	if c.TablePadding == nil {
		return ""
	}

	return *c.TablePadding
}

type TableWriter struct {
	config TableWriterConfig
	writer *tablewriter.Table
}

func NewTableWriter(output io.Writer, config TableWriterConfig) *TableWriter {
	return &TableWriter{config: config, writer: tablewriter.NewWriter(os.Stdout)}
}

func (tw TableWriter) RenderTable(rows [][]string) {
	tw.writer.SetCenterSeparator(tw.config.getCenterSeparator())
	tw.writer.SetColumnSeparator(tw.config.getColumnSeparator())
	tw.writer.SetColumnAlignment(tw.config.getColumnAlignment())
	tw.writer.SetRowSeparator(tw.config.getRowSeparator())
	tw.writer.SetHeaderLine(tw.config.getHeaderLine())
	tw.writer.SetBorder(tw.config.getBorder())
	tw.writer.SetTablePadding(tw.config.getTablePadding())

	tw.writer.SetAutoFormatHeaders(false)

	tw.writer.SetHeader(tw.config.getHeader())

	for _, row := range rows {
		tw.writer.Append(row)
	}

	tw.writer.Render()
}
