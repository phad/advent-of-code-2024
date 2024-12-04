package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const xmas = "XMAS"

var lineRE = regexp.MustCompile("mul\\(([0-9]+),([0-9]+)\\)")

type grid struct {
	w, h  int
	cells [][]rune
}

func newGrid(in []string) (*grid, error) {
	g := &grid{h: len(in)}
	for i, r := range in {
		if i == 0 {
			g.w = len(r)
			if g.w != g.h {
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

func (g *grid) Highlight(show string) string {
	s := fmt.Sprintf("width:%d height:%d\n", g.w, g.h)
	for _, row := range g.cells {
		r := make([]rune, g.w)
		copy(r, row)
		for i, c := range row {
			if !strings.ContainsRune(show, c) {
				r[i] = '.'
			}
		}
		s += string(r)
		s += "\n"
	}
	return s
}

func (g *grid) numHoriz(s string) int {
	if len(s) > g.w || len(s) == 0 {
		return 0
	}
	var check []string
	for _, r := range g.cells {
		s := string(r)
		check = append(check, s)
		check = append(check, reverse(s))
	}
	return countAll(s, check)
}

func (g *grid) numVert(s string) int {
	if len(s) > g.h || len(s) == 0 {
		return 0
	}
	var check []string
	for x := 0; x < g.w; x++ {
		var col []rune
		for y := 0; y < g.h; y++ {
			col = append(col, g.cells[y][x])
		}
		check = append(check, string(col))
		check = append(check, reverse(string(col)))
	}
	return countAll(s, check)
}

func (g *grid) numDiag1(s string) int {
	var check []string
	for y := 0; y < 2*g.h-1; y++ {
		var diag []rune
		for x := 0; x <= y; x++ {
			yy := y - x
			if yy >= g.h || x >= g.w {
				continue
			}
			//log.Printf("y:%d,x:%d", yy, x)
			diag = append(diag, g.cells[yy][x])
		}
		//log.Printf("diag1 %d: %q", y, string(diag))

		check = append(check, string(diag))
		check = append(check, reverse(string(diag)))
	}
	return countAll(s, check)
}

func (g *grid) numDiag2(s string) int {
	var check []string
	for y := 0; y < 2*g.h-1; y++ {
		var diag []rune
		for x := g.w - 1; x >= g.w-1-y; x-- {
			yy := y - (g.w - 1 - x)
			if yy >= g.h || x < 0 {
				continue
			}
			//log.Printf("y:%d,x:%d", yy, x)
			diag = append(diag, g.cells[yy][x])
		}
		//log.Printf("diag2 %d: %q", y, string(diag))

		check = append(check, string(diag))
		check = append(check, reverse(string(diag)))
	}
	return countAll(s, check)
}

func reverse(s string) string {
	l := len(s)
	if l < 2 {
		return s
	}
	r := []rune(s)
	for i := 0; i < l/2; i++ {
		r[i], r[l-i-1] = r[l-i-1], r[i]
	}
	return string(r)
}

func countAll(s string, check []string) int {
	n := 0
	for _, c := range check {
		n += strings.Count(c, s)
	}
	return n
}

func main() {
	log.Println("AoC-2024-day04-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	g, err := newGrid(lines)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Grid:\n%v", g.Highlight(xmas))

	nh := g.numHoriz(xmas)
	nv := g.numVert(xmas)
	nd1 := g.numDiag1(xmas)
	nd2 := g.numDiag2(xmas)

	log.Printf("nh:%d nv:%d nd1:%d nd2:%d", nh, nv, nd1, nd2)
	log.Printf("found %d matches", nh+nv+nd1+nd2)
}
