package graph

import (
	"github.com/idursun/jjui/internal/screen"
	"io"
	"strings"
	"unicode/utf8"
)

func ParseRows(reader io.Reader) []Row {
	var rows []Row
	var row Row
	rawSegments := screen.ParseFromReader(reader)

	for segmentedLine := range screen.BreakNewLinesIter(rawSegments) {
		rowLine := NewGraphRowLine(segmentedLine)
		if changeIdIdx := rowLine.FindPossibleChangeIdIdx(); changeIdIdx != -1 {
			rowLine.Flags = Revision | Highlightable
			previousRow := row
			row = NewGraphRow()
			if previousRow.Commit != nil {
				rows = append(rows, previousRow)
				row.Previous = &previousRow
			}
			for j := 0; j < changeIdIdx; j++ {
				row.Indent += utf8.RuneCountInString(rowLine.Segments[j].Text)
			}
			rowLine.ChangeIdIdx = changeIdIdx
			row.Commit.ChangeId = rowLine.Segments[changeIdIdx].Text
			for nextIdx := changeIdIdx + 1; nextIdx < len(rowLine.Segments); nextIdx++ {
				nextSegment := rowLine.Segments[nextIdx]
				if strings.TrimSpace(nextSegment.Text) == "" || strings.ContainsAny(nextSegment.Text, "\n\t\r ") {
					break
				}
				row.Commit.ChangeId += nextSegment.Text
			}
			if commitIdIdx := rowLine.FindPossibleCommitIdIdx(changeIdIdx); commitIdIdx != -1 {
				rowLine.CommitIdIdx = commitIdIdx
				row.Commit.CommitId = rowLine.Segments[commitIdIdx].Text
			}
		}
		row.AddLine(&rowLine)
	}
	if row.Commit != nil {
		rows = append(rows, row)
	}
	return rows
}
