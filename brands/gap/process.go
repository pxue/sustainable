package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/pxue/sustainable/ocs"
)

var (
	Debug = false
)

type pageHelper struct {
	width        float64
	height       float64
	columnBounds []float64
	tableBounds  []float64
	rowBounds    []float64
}

func logBox(box *ocs.BoundingBox, word string) {
	for _, v := range box.NormalizedVertices {
		fmt.Printf("(%.3f %.3f), ", v.X, v.Y)
	}
	fmt.Printf("%s\n", word)
}

func debugf(format string, args ...interface{}) {
	if !Debug {
		return
	}
	fmt.Printf(format, args...)
}

func (p *pageHelper) findColumn(v float64) int {
	for i := 0; i < len(p.columnBounds)-1; i++ {
		if p.toX(v) > p.columnBounds[i] && p.toX(v) < p.columnBounds[i+1] {
			return i
		}
	}
	return -1
}

func (p *pageHelper) findRow(v float64) int {
	if v > p.rowBounds[len(p.rowBounds)-1] {
		return len(p.rowBounds)
	}
	for i := 0; i < len(p.rowBounds)-1; i++ {
		// bottom bounded.
		if v > p.rowBounds[i] && v < p.rowBounds[i+1] {
			return i
		}
	}
	return -1
}

func (p *pageHelper) toX(v float64) float64 {
	return v * p.width
}

func (p *pageHelper) toY(v float64) float64 {
	return v * p.height
}

func process(w io.Writer, r io.Reader) {
	var wrapper *ocs.Wrapper
	if err := json.NewDecoder(r).Decode(&wrapper); err != nil {
		log.Fatal(err)
	}
	for _, response := range wrapper.Responses {
		if response.Context.PageNumber == 1 {
			continue
		}
		for _, page := range response.Annotation.Pages {
			p := &pageHelper{
				width:        page.Width,
				height:       page.Height,
				columnBounds: []float64{0.0, 240, 435, 505, 570, 633, 695, 792}, // hardcoded. can we do better?
				tableBounds:  []float64{95, 85, 755, 545},
				rowBounds:    []float64{0}, // bottom of the row
			}

			// intermedia rows.
			pageWords := []*ocs.Word{}
			for _, block := range page.Blocks {
				for _, para := range block.Paragraphs {
					for _, word := range para.Words {
						wv := word.BoundingBox.NormalizedVertices
						if p.toX(wv[0].X) < p.tableBounds[0] || p.toY(wv[0].Y) < p.tableBounds[1] ||
							p.toX(wv[0].X) > p.tableBounds[2] || p.toY(wv[0].Y) > p.tableBounds[3] {
							// ignore outside of table content
							continue
						}
						// got the row, append to the correct column
						idx := p.findColumn(wv[0].X)
						if idx == 5 {
							// anchor word
							b := word.BoundingBox.NormalizedVertices[3].Y
							p.rowBounds = append(p.rowBounds, b)
						}
						pageWords = append(pageWords, word)
					}
				}
			}

			sort.Sort(sort.Float64Slice(p.rowBounds))
			fmt.Println(p.rowBounds, len(pageWords))
			rows := make([]*ocs.Row, len(p.rowBounds))

			for _, w := range pageWords {
				wcy := w.BoundingBox.NormalizedVertices[3].Y - (w.Height() / 2.0)
				//fmt.Printf("\tgot word(%s,%.5f)->%d\n", w, wcy, p.findRow(wcy))
				rdx := p.findRow(wcy)

				row := rows[rdx]
				if row == nil {
					row = ocs.NewRow(rdx, 7)
					rows[rdx] = row
				}
				cdx := p.findColumn(w.BoundingBox.NormalizedVertices[0].X)
				col := row.Columns[cdx]

				col.Words = append(col.Words, w)
			}

			for i, row := range rows {
				fmt.Fprintf(w, "%d,%s\n", i, row)
			}
		}
		fmt.Printf("finished page %d\n\n", response.Context.PageNumber)

	}
}

func main() {
	i := 3
	log.Printf("parsing %d", i)

	writer, err := os.Create(fmt.Sprintf("gap%d.csv", i))
	if err != nil {
		log.Fatal(err)
	}
	//writer := os.Stdout
	reader, err := os.Open(fmt.Sprintf("gap%d.json", i))
	if err != nil {
		log.Fatal(err)
	}
	process(writer, reader)
	reader.Close()

	log.Printf("done %d", i)

	writer.Close()
}
