package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/phad/advent-of-code-2024/aoc"
)

const xmas = "MAS"

var lineRE = regexp.MustCompile("mul\\(([0-9]+),([0-9]+)\\)")

type coord struct {
	y, x int
}

type coordRune struct {
	c coord
	r rune
}

// highlight renders the grid, replacing any rune not present in show with '.'.
func highlight(g *aoc.Grid[rune], show string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "width:%d height:%d\n", g.W, g.H)
	for _, row := range g.Cells {
		r := make([]rune, g.W)
		copy(r, row)
		for i, c := range row {
			if !strings.ContainsRune(show, c) {
				r[i] = '.'
			}
		}
		b.WriteString(string(r))
		b.WriteByte('\n')
	}
	return b.String()
}

func coordsDiag1(g *aoc.Grid[rune], s string) []coord {
	var check [][]coordRune
	for y := 0; y < 2*g.H-1; y++ {
		var diag []coordRune
		for x := 0; x <= y; x++ {
			yy := y - x
			if yy >= g.H || x >= g.W {
				continue
			}
			diag = append(diag, coordRune{c: coord{y: yy, x: x}, r: g.Cells[yy][x]})
		}
		check = append(check, diag)
		if len(diag) > 1 {
			check = append(check, reverse(diag))
		}
	}
	return countAll(s, check)
}

func coordsDiag2(g *aoc.Grid[rune], s string) []coord {
	var check [][]coordRune
	for y := 0; y < 2*g.H-1; y++ {
		var diag []coordRune
		for x := g.W - 1; x >= g.W-1-y; x-- {
			yy := y - (g.W - 1 - x)
			if yy >= g.H || x < 0 {
				continue
			}
			diag = append(diag, coordRune{c: coord{y: yy, x: x}, r: g.Cells[yy][x]})
		}
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
	lines, err := aoc.ReadLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	g, err := aoc.NewRuneGrid(lines, true /*=wantSquare*/)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Grid:\n%v", highlight(g, xmas))

	d1 := coordsDiag1(g, xmas)
	d2 := coordsDiag2(g, xmas)

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
