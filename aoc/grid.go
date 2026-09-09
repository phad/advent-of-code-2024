package aoc

import (
	"fmt"
	"strings"
)

// Point is an (x, y) coordinate in a Grid. x indexes columns (width) and y
// indexes rows (height).
type Point struct{ X, Y int }

// String renders the point as "(x,y)".
func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

// Grid is a 2-D grid of cells of type T, indexed by Point. Cells is exposed for
// callers that need to iterate directly; use At/Set for bounds-checked access.
type Grid[T comparable] struct {
	W, H  int
	Cells [][]T
	// format renders a single cell for String. Constructors set a sensible
	// default; callers may override via SetFormatter.
	format func(T) string
}

// NewGrid returns a w x h grid with every cell set to zero. The default String
// formatter uses fmt's %v for each cell.
func NewGrid[T comparable](w, h int) *Grid[T] {
	g := &Grid[T]{W: w, H: h, format: func(v T) string { return fmt.Sprintf("%v", v) }}
	g.Cells = make([][]T, h)
	for y := range g.Cells {
		g.Cells[y] = make([]T, w)
	}
	return g
}

// SetFormatter overrides how String renders each cell.
func (g *Grid[T]) SetFormatter(f func(T) string) {
	g.format = f
}

// In reports whether p lies within the grid bounds.
func (g *Grid[T]) In(p Point) bool {
	return p.X >= 0 && p.X < g.W && p.Y >= 0 && p.Y < g.H
}

// At returns the cell at p, and whether p was in bounds. When out of bounds it
// returns the zero value and false.
func (g *Grid[T]) At(p Point) (T, bool) {
	if !g.In(p) {
		var zero T
		return zero, false
	}
	return g.Cells[p.Y][p.X], true
}

// Set writes v at p, returning whether p was in bounds.
func (g *Grid[T]) Set(p Point, v T) bool {
	if !g.In(p) {
		return false
	}
	g.Cells[p.Y][p.X] = v
	return true
}

// Swap exchanges the cells at p1 and p2, returning whether both were in bounds.
func (g *Grid[T]) Swap(p1, p2 Point) bool {
	if !g.In(p1) || !g.In(p2) {
		return false
	}
	g.Cells[p1.Y][p1.X], g.Cells[p2.Y][p2.X] = g.Cells[p2.Y][p2.X], g.Cells[p1.Y][p1.X]
	return true
}

// Find returns the first cell (scanning rows top-to-bottom, left-to-right)
// equal to v, and whether one was found.
func (g *Grid[T]) Find(v T) (Point, bool) {
	var found Point
	var ok bool
	g.FindAll(v, func(p Point) bool {
		found, ok = p, true
		return false
	})
	return found, ok
}

// FindAll calls f for each cell equal to v. If f returns false the scan stops.
func (g *Grid[T]) FindAll(v T, f func(Point) bool) {
	for y, row := range g.Cells {
		for x, c := range row {
			if c == v && !f(Point{x, y}) {
				return
			}
		}
	}
}

// Count returns the number of cells equal to v.
func (g *Grid[T]) Count(v T) int {
	n := 0
	for _, row := range g.Cells {
		for _, c := range row {
			if c == v {
				n++
			}
		}
	}
	return n
}

// String renders the grid with a "width:W height:H" header followed by one
// line per row, using the instance formatter for each cell.
func (g *Grid[T]) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "width:%d height:%d\n", g.W, g.H)
	for _, row := range g.Cells {
		for _, c := range row {
			b.WriteString(g.format(c))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// checkDims validates that in is rectangular (and optionally square), returning
// its width and height.
func checkDims(in []string, wantSquare bool) (w, h int, err error) {
	h = len(in)
	for i, r := range in {
		if i == 0 {
			w = len(r)
			if wantSquare && w != h {
				return 0, 0, fmt.Errorf("grid isn't square: width %d != height %d", w, h)
			}
		}
		if len(r) != w {
			return 0, 0, fmt.Errorf("row %d wrong size %d want %d", i, len(r), w)
		}
	}
	return w, h, nil
}

// NewRuneGrid builds a Grid[rune] from input lines, one rune per column. When
// wantSquare is true it errors unless width == height. String renders each cell
// as its rune.
func NewRuneGrid(in []string, wantSquare bool) (*Grid[rune], error) {
	w, h, err := checkDims(in, wantSquare)
	if err != nil {
		return nil, err
	}
	g := &Grid[rune]{W: w, H: h, format: func(r rune) string { return string(r) }}
	for _, line := range in {
		g.Cells = append(g.Cells, []rune(line))
	}
	return g, nil
}

// NewDigitGrid builds a Grid[int] from input lines, parsing each character as a
// single base-10 digit. When wantSquare is true it errors unless width ==
// height. String renders each cell as its digit.
func NewDigitGrid(in []string, wantSquare bool) (*Grid[int], error) {
	w, h, err := checkDims(in, wantSquare)
	if err != nil {
		return nil, err
	}
	g := &Grid[int]{W: w, H: h, format: func(d int) string { return fmt.Sprintf("%d", d) }}
	for _, line := range in {
		row := make([]int, w)
		for x := 0; x < w; x++ {
			row[x] = int(MustParseInt(line[x : x+1]))
		}
		g.Cells = append(g.Cells, row)
	}
	return g, nil
}

// Highlight renders the grid like String but replaces every cell not equal to
// show with '.'. It is defined for rune grids, matching the common
// per-day helper.
func (g *Grid[T]) Highlight(show T) string {
	var b strings.Builder
	fmt.Fprintf(&b, "width:%d height:%d\n", g.W, g.H)
	dot := "."
	for _, row := range g.Cells {
		for _, c := range row {
			if c == show {
				b.WriteString(g.format(c))
			} else {
				b.WriteString(dot)
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}
