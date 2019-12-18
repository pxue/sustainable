package ocs

import "fmt"

type Vertice struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type BoundingBox struct {
	NormalizedVertices []*Vertice `json:"normalizedVertices"`
}

type Break struct {
	Type string `json:"type"`
}

type Symbol struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
	Property   struct {
		DetectedBreak *Break `json:"detectedBreak"`
	} `json:"property,omitempty"`
}

type Word struct {
	BoundingBox *BoundingBox `json:"boundingBox"`
	Symbols     []*Symbol    `json:"symbols"`
	Confidence  float64      `json:"confidence"`

	Hidden bool `json:"-"`
}

const (
	LineBreak = "LINE_BREAK"
	Space     = "SPACE"
	EOLSpace  = "EOL_SURE_SPACE"
)

func (w *Word) String() string {
	str := ""
	for _, sym := range w.Symbols {
		str += sym.Text
	}
	return str
}

func (w *Word) DebugStr() string {
	str := ""
	for _, sym := range w.Symbols {
		str += sym.Text
		if b := sym.Property.DetectedBreak; b != nil {
			str = fmt.Sprintf("%s (%s)", str, b.Type)
		}
	}
	for _, v := range w.BoundingBox.NormalizedVertices {
		str = fmt.Sprintf("%s (%.8f,%.8f)", str, v.X, v.Y)
	}
	return str
}

func (w *Word) GetBreak() *Break {
	sym := w.Symbols[len(w.Symbols)-1]
	return sym.Property.DetectedBreak
}

func (w *Word) Height() float64 {
	top := w.BoundingBox.NormalizedVertices[0].Y
	bot := w.BoundingBox.NormalizedVertices[2].Y
	return bot - top
}

func (w *Word) Width() float64 {
	left := w.BoundingBox.NormalizedVertices[0].X
	right := w.BoundingBox.NormalizedVertices[1].X
	return right - left
}

type Paragraph struct {
	BoundingBox *BoundingBox `json:"boundingBox"`
	Words       []*Word      `json:"words"`
	Confidence  float64      `json:"confidence"`
}

type Block struct {
	BoundingBox *BoundingBox `json:"boundingBox"`
	Paragraphs  []*Paragraph `json:"paragraphs"`
	BlockType   string       `json:"blockType"`
	Confidence  float64      `json:"confidence"`
}

type Page struct {
	Property struct {
		DetectedLanguages []struct {
			LanguageCode string  `json:"languageCode"`
			Confidence   float64 `json:"confidence"`
		} `json:"detectedLanguages"`
	} `json:"property"`
	Width  float64  `json:"width"`
	Height float64  `json:"height"`
	Blocks []*Block `json:"blocks"`
}

type FullTextAnnotation struct {
	Pages []*Page `json:"pages"`
	Text  string  `json:"text"`
}

type Context struct {
	URI        string `json:"uri"`
	PageNumber int    `json:"pageNumber"`
}

type Response struct {
	Annotation *FullTextAnnotation `json:"fullTextAnnotation"`
	Context    *Context            `json:"context"`
}

type Wrapper struct {
	Responses []*Response `json:"responses"`
}
