package ocs

type ByX []*Vertice

func (a ByX) Len() int           { return len(a) }
func (a ByX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByX) Less(i, j int) bool { return a[i].X < a[j].X }

type ByY []*Vertice

func (a ByY) Len() int           { return len(a) }
func (a ByY) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByY) Less(i, j int) bool { return a[i].Y < a[j].Y }

type SortedWords []*Word

func (s SortedWords) Len() int      { return len(s) }
func (s SortedWords) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortedWords) Less(i, j int) bool {
	siv := s[i].BoundingBox.NormalizedVertices
	sjv := s[j].BoundingBox.NormalizedVertices
	if siv[0].Y < sjv[0].Y {
		// different lines
		return true
	} else {
		if siv[0].Y == sjv[0].Y {
			// same line.
			if siv[0].X < sjv[0].X {
				return true
			}
		}
	}
	return false
}
