package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const xmas = "MAS"

var lineRE = regexp.MustCompile("mul\\(([0-9]+),([0-9]+)\\)")

type grid struct {
	w, h  int
	cells [][]rune
}

type coord struct {
	y, x int
}

type coordRune struct {
	c coord
	r rune
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

func (g *grid) highlight(show string) string {
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

func (g *grid) coordsDiag1(s string) []coord {
	var check [][]coordRune
	for y := 0; y < 2*g.h-1; y++ {
		var diag []coordRune
		for x := 0; x <= y; x++ {
			yy := y - x
			if yy >= g.h || x >= g.w {
				continue
			}
			//log.Printf("y:%d,x:%d", yy, x)
			diag = append(diag, coordRune{c: coord{y: yy, x: x}, r: g.cells[yy][x]})
		}
		//log.Printf("diag1 %d: %q", y, string(diag))

		check = append(check, diag)
		if len(diag) > 1 {
			check = append(check, reverse(diag))
		}
	}
	return countAll(s, check)
}

func (g *grid) coordsDiag2(s string) []coord {
	var check [][]coordRune
	for y := 0; y < 2*g.h-1; y++ {
		var diag []coordRune
		for x := g.w - 1; x >= g.w-1-y; x-- {
			yy := y - (g.w - 1 - x)
			if yy >= g.h || x < 0 {
				continue
			}
			//log.Printf("y:%d,x:%d", yy, x)
			diag = append(diag, coordRune{c: coord{y: yy, x: x}, r: g.cells[yy][x]})
		}
		//log.Printf("diag2 %d: %q", y, string(diag))

		check = append(check, diag)
		if len(diag) > 1 {
			check = append(check, reverse(diag))
		}
	}
	return countAll(s, check)
}

func reverse(s []coordRune) []coordRune {
	l := len(s)
	if l < 2 {
		return s
	}
	r := make([]coordRune, l)
	copy(r, s)
	for i := 0; i < l/2; i++ {
		r[i], r[l-i-1] = r[l-i-1], r[i]
	}
	return r
}

func countAll(s string, check [][]coordRune) []coord {
	var found []coord
	for _, crs := range check {
		//log.Printf("check #%d: %v", i, crs)
		var rs []rune
		for _, cr := range crs {
			rs = append(rs, cr.r)
		}
		//log.Printf("  -> %v", string(rs))
		for k := 0; k <= len(rs)-len(s); k++ {
			str := string(rs[k:])
			idx := strings.Index(str, s)
			if idx == -1 {
				continue
			}
			cr := crs[k+idx+(len(s)-1)/2]
			//log.Printf("Found %v at %d: %s in %v", cr, idx, s, str)
			found = append(found, cr.c)
			k += (idx + len(s) - 1)
		}
	}
	return found
}

func main() {
	log.Println("AoC-2024-day03-part1")
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

	log.Printf("Grid:\n%v", g.highlight(xmas))

	d1 := g.coordsDiag1(xmas)
	d2 := g.coordsDiag2(xmas)

	//log.Printf("nd1:%d nd2:%d", len(d1), len(d2))
	//log.Printf("d1:%v\nd2:%v", d1, d2)

	m := map[coord]int{}
	for _, c := range d1 {
		m[c] += 1
	}
	for _, c := range d2 {
		m[c] += 1
	}

	found := 0
	for _, num := range m {
		if num == 2 {
			found++
		}
	}

	log.Printf("found %d X-MAS", found)
}
