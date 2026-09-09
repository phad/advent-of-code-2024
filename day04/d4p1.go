package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/phad/advent-of-code-2024/aoc"
)

const xmas = "XMAS"

var lineRE = regexp.MustCompile("mul\\(([0-9]+),([0-9]+)\\)")

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

func numHoriz(g *aoc.Grid[rune], s string) int {
	if len(s) > g.W || len(s) == 0 {
		return 0
	}
	var check []string
	for _, r := range g.Cells {
		s := string(r)
		check = append(check, s)
		check = append(check, reverse(s))
	}
	return countAll(s, check)
}

func numVert(g *aoc.Grid[rune], s string) int {
	if len(s) > g.H || len(s) == 0 {
		return 0
	}
	var check []string
	for x := 0; x < g.W; x++ {
		var col []rune
		for y := 0; y < g.H; y++ {
			col = append(col, g.Cells[y][x])
		}
		check = append(check, string(col))
		check = append(check, reverse(string(col)))
	}
	return countAll(s, check)
}

func numDiag1(g *aoc.Grid[rune], s string) int {
	var check []string
	for y := 0; y < 2*g.H-1; y++ {
		var diag []rune
		for x := 0; x <= y; x++ {
			yy := y - x
			if yy >= g.H || x >= g.W {
				continue
			}
			diag = append(diag, g.Cells[yy][x])
		}
		check = append(check, string(diag))
		check = append(check, reverse(string(diag)))
	}
	return countAll(s, check)
}

func numDiag2(g *aoc.Grid[rune], s string) int {
	var check []string
	for y := 0; y < 2*g.H-1; y++ {
		var diag []rune
		for x := g.W - 1; x >= g.W-1-y; x-- {
			yy := y - (g.W - 1 - x)
			if yy >= g.H || x < 0 {
				continue
			}
			diag = append(diag, g.Cells[yy][x])
		}
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
	lines, err := aoc.ReadLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	g, err := aoc.NewRuneGrid(lines, true /*=wantSquare*/)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Grid:\n%v", highlight(g, xmas))

	nh := numHoriz(g, xmas)
	nv := numVert(g, xmas)
	nd1 := numDiag1(g, xmas)
	nd2 := numDiag2(g, xmas)

	log.Printf("nh:%d nv:%d nd1:%d nd2:%d", nh, nv, nd1, nd2)
	log.Printf("found %d matches", nh+nv+nd1+nd2)
}
