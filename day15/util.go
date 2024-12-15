package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func readLines(f string) ([]string, error) {
	rd, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer rd.Close()

	var lines []string
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		t := scanner.Text()
		lines = append(lines, t)
	}
	return lines, scanner.Err()
}

func mustParseInt(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("ParseInt(%q) err=%v", s, err)
	}
	return v
}

type grid struct {
	w, h  int
	cells [][]rune
}

func newGrid(in []string, wantSquare bool) (*grid, error) {
	g := &grid{h: len(in)}
	for i, r := range in {
		if i == 0 {
			g.w = len(r)
			if g.w != g.h && wantSquare {
				return nil, fmt.Errorf("Grid isn't square: width %d != height %d", g.w, g.h)
			}
		}
		if i > 0 && len(r) != g.w {
			return nil, fmt.Errorf("Row %d wrong size %d want %d", 1, len(r), g.w)
		}
		row := []rune(r)
		g.cells = append(g.cells, row)
	}
	return g, nil
}

func (g *grid) String() string {
	s := fmt.Sprintf("width:%d height:%d\n", g.w, g.h)
	for _, r := range g.cells {
		s += string(r)
		s += "\n"
	}
	return s
}

func (g *grid) at(p point) (rune, bool) {
	nope := rune(0)
	if p.y < 0 || p.y >= g.h || p.x < 0 || p.y >= g.w {
		return nope, false
	}
	return g.cells[p.y][p.x], true
}

func (g *grid) find(r rune) (point, bool) {
	var found point
	var ok bool
	g.findAll(r, func(p point) bool {
		found, ok = p, true
		return false
	})
	return found, ok
}

// f is called for each find. If f returns false the search stops.
func (g *grid) findAll(r rune, f func(point) bool) {
	for y, row := range g.cells {
		for x, c := range row {
			if r == c && !f(point{x, y}) {
				return
			}
		}
	}
}

func (g *grid) swap(p1, p2 point) bool {
	if p1.y < 0 || p1.y >= g.h || p1.x < 0 || p1.y >= g.w {
		return false
	}
	if p2.y < 0 || p2.y >= g.h || p2.x < 0 || p2.y >= g.w {
		return false
	}
	g.cells[p2.y][p2.x], g.cells[p1.y][p1.x] = g.cells[p1.y][p1.x], g.cells[p2.y][p2.x]
	return true
}

func (g *grid) highlight(show rune) string {
	s := fmt.Sprintf("width:%d height:%d\n", g.w, g.h)
	for _, row := range g.cells {
		r := make([]rune, g.w)
		copy(r, row)
		for i, c := range row {
			if show != c {
				r[i] = '.'
			}
		}
		s += string(r)
		s += "\n"
	}
	return s
}

type point struct{ x, y int }
