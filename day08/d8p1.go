package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/phad/advent-of-code-2024/aoc"
)

/* input format
............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............
*/

// set marks p in g. Out-of-bounds points are ignored.
func set(g *aoc.Grid[bool], p aoc.Point) {
	g.Set(p, true)
}

// unset clears p in g. Out-of-bounds points are ignored.
func unset(g *aoc.Grid[bool], p aoc.Point) {
	g.Set(p, false)
}

// get reports whether p is marked in g.
func get(g *aoc.Grid[bool], p aoc.Point) bool {
	v, _ := g.At(p)
	return v
}

// numSet returns the number of marked cells in g.
func numSet(g *aoc.Grid[bool]) int {
	return g.Count(true)
}

func union(g, other *aoc.Grid[bool]) error {
	if g.W != other.W || g.H != other.H {
		return fmt.Errorf("mismatched grid sizes!")
	}
	for y := 0; y < g.H; y++ {
		for x := 0; x < g.W; x++ {
			p := aoc.Point{X: x, Y: y}
			if get(other, p) {
				set(g, p)
			}
		}
	}
	return nil
}

func remove(g, other *aoc.Grid[bool]) error {
	if g.W != other.W || g.H != other.H {
		return fmt.Errorf("mismatched grid sizes!")
	}
	for y := 0; y < g.H; y++ {
		for x := 0; x < g.W; x++ {
			p := aoc.Point{X: x, Y: y}
			if get(other, p) {
				unset(g, p)
			}
		}
	}
	return nil
}

func assert(b bool) {
	if !b {
		log.Fatalf("boom")
	}
}

func runTest() {
	g := aoc.NewGrid[bool](2, 2)
	assert(numSet(g) == 0)
	set(g, aoc.Point{0, 0})
	assert(numSet(g) == 1)
	set(g, aoc.Point{1, 1})
	assert(numSet(g) == 2)
	unset(g, aoc.Point{0, 1})
	assert(numSet(g) == 2)
	unset(g, aoc.Point{0, 0})
	assert(numSet(g) == 1)
	unset(g, aoc.Point{0, 0})
	assert(numSet(g) == 1)
	unset(g, aoc.Point{1, 1})
	assert(numSet(g) == 0)

	set(g, aoc.Point{0, 0})
	assert(numSet(g) == 1)
	g1 := aoc.NewGrid[bool](2, 2)
	set(g1, aoc.Point{1, 1})
	assert(numSet(g1) == 1)
	union(g, g1)
	assert(numSet(g) == 2)
	remove(g1, g)
	assert(numSet(g) == 2)
	assert(numSet(g1) == 0)
	remove(g, g1)
	assert(numSet(g) == 2)
	assert(numSet(g1) == 0)
}

type antennaSet struct {
	frequency        rune
	locations        []aoc.Point
	nodes, antinodes *aoc.Grid[bool]
}

func newAntennaSet(w, h int, freq rune) *antennaSet {
	return &antennaSet{
		frequency: freq,
		nodes:     aoc.NewGrid[bool](w, h),
		antinodes: aoc.NewGrid[bool](w, h),
	}
}

func (as *antennaSet) addLocation(x, y int) {
	p := aoc.Point{X: x, Y: y}
	as.locations = append(as.locations, p)
	set(as.nodes, p)
}

func (as *antennaSet) findAntinodes() {
	for i, locN1 := range as.locations {
		for j := i + 1; j < len(as.locations); j++ {
			locN2 := as.locations[j]
			dx := locN1.X - locN2.X
			dy := locN1.Y - locN2.Y
			locAN1 := aoc.Point{
				X: locN1.X + dx,
				Y: locN1.Y + dy,
			}
			locAN2 := aoc.Point{
				X: locN2.X - dx,
				Y: locN2.Y - dy,
			}
			set(as.antinodes, locAN1)
			set(as.antinodes, locAN2)
		}
	}
}

func (as *antennaSet) String() string {
	var b strings.Builder
	for y := 0; y < as.nodes.H; y++ {
		for x := 0; x < as.nodes.W; x++ {
			p := aoc.Point{X: x, Y: y}
			if get(as.nodes, p) {
				b.WriteRune(as.frequency)
			} else if get(as.antinodes, p) {
				b.WriteRune('#')
			} else {
				b.WriteRune('.')
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func main() {
	runTest()
	log.Println("AoC-2024-day08-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := aoc.ReadLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	w, h := 0, len(lines)
	allAntennas := map[rune]*antennaSet{}

	for y, line := range lines {
		if w == 0 {
			w = len(line)
		} else if len(line) != w {
			log.Fatalf("Row %d: inconsistent row length %d want %d", y, len(line), w)
		}
		for x, ss := range strings.Split(line, "") {
			r := rune(ss[0])
			if r == '.' {
				continue
			}
			as, ok := allAntennas[r]
			if !ok {
				as = newAntennaSet(w, h, r)
				allAntennas[r] = as
			}
			as.addLocation(x, y)
		}
	}

	allNs, allANs := aoc.NewGrid[bool](w, h), aoc.NewGrid[bool](w, h)
	totalNs, totalANs := 0, 0
	for r, as := range allAntennas {
		as.findAntinodes()
		log.Printf("%v\n%v", r, as)
		union(allNs, as.nodes)
		union(allANs, as.antinodes)
		totalNs += numSet(as.nodes)
		totalANs += numSet(as.antinodes)
	}

	log.Printf("Total #nodes: %d", totalNs)
	log.Printf("Total #antinodes: %d", totalANs)
	log.Printf("Total unique #nodes: %d", numSet(allNs))
	log.Printf("Total unique #antinodes: %d <-- submit this", numSet(allANs))
}
