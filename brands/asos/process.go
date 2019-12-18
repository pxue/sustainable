package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pxue/sustainable/ocs"
)

var (
	Debug        = false
	columnBounds = []float64{0.0, 0.1871, 0.56544, 0.65, 0.73, 0.805, 0.88, 1}
)

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

func findColumn(v float64) int {
	for i := 0; i < len(columnBounds)-1; i++ {
		if v > columnBounds[i] && v < columnBounds[i+1] {
			return i
		}
	}
	return -1
}

// check the line spacing between the last words in the column
// and the next word.. if the spacing is more than 5px, it's probably
// a new row.
func shouldNewRow(col *ocs.Column, nextWord *ocs.Word) bool {
	lastWords := col.Words
	// if the column is empty. return early
	if len(lastWords) == 0 {
		return false
	}

	if col.CountLines() == 3 {
		return true
	}

	nv := nextWord.BoundingBox.NormalizedVertices
	for _, word := range col.Words {
		// for every word in the column. check if the next word
		// should go beside it
		wv := word.BoundingBox.NormalizedVertices

		// first, check if the topY values are equal
		// if we've found similar, return false
		if wv[0].Y == nv[0].Y {
			return false
		}
	}

	lastWord := lastWords[len(lastWords)-1]
	lastBotY := lastWord.BoundingBox.NormalizedVertices[3].Y
	// find the last word that's not alpha numeric
	wordY := nextWord.BoundingBox.NormalizedVertices[0].Y

	// check line spacing height
	if (wordY - lastBotY) > 0.0084 { // 5px
		return true
	}

	//check last word for breaks and check it's character
	if lastSym := lastWord.Symbols[len(lastWord.Symbols)-1]; lastSym != nil {
		if b := lastSym.Property.DetectedBreak; b != nil && !strings.Contains("(", lastSym.Text) {
			// find break. or, if last symbol is an open bracket. skip this
			// check.
			if b.Type == ocs.EOLSpace || b.Type == ocs.LineBreak {
				lastX := lastWord.BoundingBox.NormalizedVertices[1].X
				col := findColumn(lastX)
				boundB := columnBounds[col+1]

				// check if the nextWord could have fit beside lastWord
				debugf("\t\tcould we have fit on the same line? lw(%s), (%.5f:%.5f)->%v\n", lastWord, lastX+nextWord.Width(), boundB, boundB > lastX+nextWord.Width())
				//if lastX+nextWord.Width() < boundB {
				//return true
				//}

				if lastX < 0.8*boundB {
					return true
				}
			}
		}
	}

	return false
}

func process(w io.Writer, r io.Reader) error {
	var wrapper *ocs.Wrapper
	if err := json.NewDecoder(r).Decode(&wrapper); err != nil {
		log.Fatal(err)
	}
	for _, response := range wrapper.Responses {
		//if response.Context.PageNumber != 15 {
		//continue
		//}
		for _, page := range response.Annotation.Pages {
			row := ocs.NewRow(0)
			rowsLookup := map[int]*ocs.Row{
				0: row,
			}
			for _, block := range page.Blocks {
				debugf("\tnew block\n")
				//for row.IsFilled() {
				////debugf("\trow(%d) is filled\n", row.Index)

				//nextRow, found := rowsLookup[row.Index+1]
				//if found {
				//// continue, check again
				//row = nextRow
				////debugf("\tchecking next row(%d)\n", row.Index)
				//continue
				//}
				//row = NewRow(row.Index + 1)
				//rowsLookup[row.Index] = row

				////debugf("\tswitching to new row(%d)\n", row.Index)
				//}
				//debugf("\tcurrent row(%d): %s\n\n", row.Index, row)
				for _, para := range block.Paragraphs {
					for _, word := range para.Words {
						// skip the first row
						if word.BoundingBox.NormalizedVertices[2].Y <= 0.1 {
							continue
						}
						debugf("\n\t\tgot word: %s(%d)\n", word.DebugStr(), len(word.String()))

						if len(word.String()) == 1 && strings.Contains("|", string(word.String()[0])) {
							continue
						}

						idx := findColumn(word.BoundingBox.NormalizedVertices[0].X)
						col := row.Columns[idx]

						nextRow := row
						nextCol := col
						for shouldNewRow(nextCol, word) {
							debugf("\t\t\ttested row(%d): %s\n", nextRow.Index, nextCol)
							// figure out where the word should go depending on if
							// the column alread has 2 lines or next word is more
							// spaced out than needed
							next, found := rowsLookup[nextRow.Index+1]
							if found {
								nextRow = next
								nextCol = nextRow.Columns[idx]
								continue
							}
							nextRow = ocs.NewRow(nextRow.Index + 1)
							nextCol = nextRow.Columns[idx]
							rowsLookup[nextRow.Index] = nextRow
						}
						nextCol.Words = append(nextCol.Words, word)
						debugf("\t\tcurrent row(%d)->ended up on row(%d): %s\n", row.Index, nextRow.Index, nextRow)
					}
					//for _, v := range para.BoundingBox.NormalizedVertices {
					//debugf("(%.3f %.3f), ", v.X*float64(page.Width), v.Y*float64(page.Height))
					//}
					//debugf("\n idx: %d -> %s\n\n", row.Index, row)
				}
			}

			rows := make([]*ocs.Row, len(rowsLookup))
			for idx, row := range rowsLookup {
				rows[idx] = row
			}

			for i, row := range rows {
				fmt.Fprintf(w, "%d,%s\n", i, row)
			}
		}

	}
	return nil
}

func main() {
	writer, err := os.Create("asos.csv")
	if err != nil {
		log.Fatal(err)
	}
	//writer := os.Stdout

	i := 3
	log.Printf("parsing %d", i)
	reader, err := os.Open(fmt.Sprintf("asos%d.json", i))
	if err != nil {
		log.Fatal(err)
	}
	process(writer, reader)
	reader.Close()

	log.Printf("done %d", i)

	writer.Close()
}
