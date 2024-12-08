package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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

type grid struct {
	w, h   int
	points [][]bool
	num    int
}

func newGrid(w, h int) *grid {
	g := &grid{w: w, h: h}
	for i := 0; i < h; i++ {
		g.points = append(g.points, make([]bool, w))
	}
	return g
}

type point struct{ x, y int }

func (g *grid) outside(p point) bool {
	return p.x < 0 || p.x >= g.w || p.y < 0 || p.y >= g.h
}

func (g *grid) set(p point) {
	if g.outside(p) {
		return
	}
	if g.points[p.y][p.x] {
		return
	}
	g.points[p.y][p.x] = true
	g.num++
}

func (g *grid) unset(p point) {
	if g.outside(p) {
		return
	}
	if !g.points[p.y][p.x] {
		return
	}
	g.points[p.y][p.x] = false
	g.num--
}

func (g *grid) get(p point) bool {
	return g.points[p.y][p.x]
}

func (g *grid) numSet() int {
	return g.num
}

func (g *grid) union(other *grid) error {
	if g.w != other.w || g.h != other.h {
		return fmt.Errorf("mismatched grid sizes!")
	}
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			p := point{x: x, y: y}
			if other.get(p) {
				g.set(p)
			}
		}
	}
	return nil
}

func (g *grid) remove(other *grid) error {
	if g.w != other.w || g.h != other.h {
		return fmt.Errorf("mismatched grid sizes!")
	}
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			p := point{x: x, y: y}
			if other.get(p) {
				g.unset(p)
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
	g := newGrid(2, 2)
	assert(g.numSet() == 0)
	g.set(point{0, 0})
	assert(g.numSet() == 1)
	g.set(point{1, 1})
	assert(g.numSet() == 2)
	g.unset(point{0, 1})
	assert(g.numSet() == 2)
	g.unset(point{0, 0})
	assert(g.numSet() == 1)
	g.unset(point{0, 0})
	assert(g.numSet() == 1)
	g.unset(point{1, 1})
	assert(g.numSet() == 0)

	g.set(point{0, 0})
	assert(g.numSet() == 1)
	g1 := newGrid(2, 2)
	g1.set(point{1, 1})
	assert(g1.numSet() == 1)
	g.union(g1)
	assert(g.numSet() == 2)
	g1.remove(g)
	assert(g.numSet() == 2)
	assert(g1.numSet() == 0)
	g.remove(g1)
	assert(g.numSet() == 2)
	assert(g1.numSet() == 0)
}

type antennaSet struct {
	frequency        rune
	locations        []point
	nodes, antinodes *grid
}

func newAntennaSet(w, h int, freq rune) *antennaSet {
	return &antennaSet{
		frequency: freq,
		nodes:     newGrid(w, h),
		antinodes: newGrid(w, h),
	}
}

func (as *antennaSet) addLocation(x, y int) {
	p := point{x: x, y: y}
	as.locations = append(as.locations, p)
	as.nodes.set(p)
}

func (as *antennaSet) findAntinodes() {
	for i, locN1 := range as.locations {
		for j := i + 1; j < len(as.locations); j++ {
			locN2 := as.locations[j]

			dx := locN1.x - locN2.x
			dy := locN1.y - locN2.y

			locAN1, locAN2 := locN1, locN2
			for {
				as.antinodes.set(locAN1)
				as.antinodes.set(locAN2)
				locAN1.x += dx
				locAN1.y += dy
				locAN2.x -= dx
				locAN2.y -= dy

				if as.antinodes.outside(locAN1) && as.antinodes.outside(locAN2) {
					break
				}
			}
		}
	}
}

func (as *antennaSet) String() string {
	var b strings.Builder
	for y := 0; y < as.nodes.h; y++ {
		for x := 0; x < as.nodes.w; x++ {
			p := point{x: x, y: y}
			if as.nodes.get(p) {
				b.WriteRune(as.frequency)
			} else if as.antinodes.get(p) {
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
	log.Println("AoC-2024-day08-part2")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
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

	allNs, allANs := newGrid(w, h), newGrid(w, h)
	totalNs, totalANs := 0, 0
	for r, as := range allAntennas {
		as.findAntinodes()
		log.Printf("%v\n%v", r, as)
		allNs.union(as.nodes)
		allANs.union(as.antinodes)
		totalNs += as.nodes.numSet()
		totalANs += as.antinodes.numSet()
	}

	log.Printf("Total #nodes: %d", totalNs)
	log.Printf("Total #antinodes: %d", totalANs)
	log.Printf("Total unique #nodes: %d", allNs.numSet())
	log.Printf("Total unique #antinodes: %d <-- submit this", allANs.numSet())
}
