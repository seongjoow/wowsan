package utils

import (
	"fmt"
	"strings"
)

type Table struct {
	headers []string
	lengths []int
	rows    [][]string
	title   string
}

func NewTable(headers []string, lengths []int) *Table {

	return &Table{
		headers: headers,
		lengths: lengths,
	}
}

func NewFromToTable() *Table {
	headers := []string{"From", "To"}
	lengths := []int{15, 15}
	return NewTable(headers, lengths)
}

func (t *Table) SetTitle(title string) {
	t.title = title
}

func (t *Table) AddRow(row []string) {
	if len(row) != len(t.lengths) {
		fmt.Println("Error: The number of elements in the row must match the number of column lengths")
		return
	}
	t.rows = append(t.rows, row)
}

/*
+-----------------------------------------------------------------------+
|                                 Title                                 | optional
+-----------------+-----------------+-----------------+-----------------+
|    Header1      |    Header2      |   Header3       |    Header4      |
+-----------------+-----------------+-----------------+-----------------+
|    Row1.Data1   |    Row1.Data2   |   Row1.Data3    |    Row1.Data4   |
|-----------------|-----------------|-----------------|-----------------|
|    Row2.Data1   |    Row2.Data2   |   Row2.Data3    |    Row2.Data4   |
|-----------------|-----------------|-----------------|-----------------|
*/
func (t *Table) PrintTable() {
	t.PrintTableTitle()
	t.PrintSeparatorLine()
	t.PrintHeader()
	t.PrintSeparatorLine()

	// Print each row of data
	for _, row := range t.rows {
		fmt.Printf(t.generateDataFormat(), t.formatRow(row)...)
	}

	// Print the separator line at the end
	t.PrintSeparatorLine()
}

func (t *Table) PrintTableTitle() {
	if t.title != "" {
		var titleLine string

		totalWidth := 0
		for _, length := range t.lengths {
			totalWidth += length + 1
		}
		totalWidth-- // Adjust for the last column border
		titleLine = "+" + strings.Repeat("-", totalWidth) + "+\n|" + CenterString(t.title, totalWidth) + "|"
		fmt.Println(titleLine)
	}
}

func (t *Table) PrintHeader() {
	var headerLine string
	for i, header := range t.headers {
		headerLine += "|" + CenterString(header, t.lengths[i])
	}
	headerLine += "|"
	fmt.Println(headerLine)
}

func (t *Table) PrintSeparatorLine() {
	var separatorLine string
	for _, length := range t.lengths {
		separatorLine += "+" + strings.Repeat("-", length)
	}
	separatorLine += "+"
	fmt.Println(separatorLine)
}

func (t *Table) PrintRow(row []string) {
	fmt.Printf(t.generateDataFormat(), t.formatRow(row)...)
}

func (t *Table) generateDataFormat() string {
	var formatString string
	for _, length := range t.lengths {
		formatString += "|%-" + fmt.Sprintf("%ds", length)
	}
	formatString += "|\n"
	return formatString
}

func (t *Table) formatRow(row []string) []interface{} {
	formattedRow := make([]interface{}, len(row))
	for i, col := range row {
		ellipsis := "..."
		keepLength := t.lengths[i] - len(ellipsis)
		prefixChars := keepLength / 2
		suffixChars := keepLength - prefixChars

		formattedRow[i] = CenterString(TrimString(col, t.lengths[i], prefixChars, suffixChars, ellipsis), t.lengths[i])
	}
	return formattedRow
}
