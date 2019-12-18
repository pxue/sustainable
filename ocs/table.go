package ocs

import (
	"fmt"
	"sort"
	"strings"
)

type Table struct {
	Rows []*Row
}

type Column struct {
	BoundingBox *BoundingBox
	Index       int
	Words       []*Word
}

func (c *Column) String() string {
	words := []string{}
	sort.Sort(SortedWords(c.Words))
	for _, word := range c.Words {
		words = append(words, word.String())
	}
	return strings.TrimSpace(strings.Join(words, " "))
}

func (c *Column) Height() float64 {
	for _, w := range c.Words {
		return w.Height()
	}
	return 0
}

func (c *Column) CountLines() int {
	lines := 0
	for _, w := range c.Words {
		if b := w.GetBreak(); b != nil && b.Type != Space {
			lines++
		}
	}
	return lines
}

type Row struct {
	BoundingBox *BoundingBox
	Columns     []*Column
	Index       int
}

func (r *Row) String() string {
	rowText := []string{}
	for _, c := range r.Columns {
		rowText = append(rowText, fmt.Sprintf(`"%s"`, c.String()))
	}
	return strings.Join(rowText, ",")
}

func (r *Row) Height() float64 {
	for _, c := range r.Columns {
		return c.Height()
	}
	return 0
}

func NewRow(i, size int) *Row {
	row := &Row{
		Columns: make([]*Column, size),
		Index:   i,
	}
	for i := range row.Columns {
		row.Columns[i] = &Column{
			Index: i + 1,
			Words: []*Word{},
		}
	}
	return row
}

func (r *Row) IsFilled() bool {
	for _, c := range r.Columns {
		if len(c.Words) == 0 {
			return false
		}
	}
	return true
}
