package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type grid struct {
	w, h  int
	cells [][]int
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
		var row []int
		for j := 0; j < g.w; j++ {
			row = append(row, int(mustParseInt(r[j:j+1])))
		}
		g.cells = append(g.cells, row)
	}
	return g, nil
}

func (g *grid) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("width:%d height:%d\n", g.w, g.h))
	for _, r := range g.cells {
		for _, c := range r {
			s.WriteString(fmt.Sprintf("%d", c))
		}
		s.WriteRune('\n')
	}
	return s.String()
}

type point struct{ x, y int }

func (p point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

func (g *grid) heightAt(p point) int {
	if p.x < 0 || p.x >= g.w || p.y < 0 || p.y >= g.h {
		log.Fatalf("Point %v is outside the grid!", p)
	}
	return g.cells[p.y][p.x]
}

type route []point

type trailhead struct {
	start  point
	routes []route
}

func (th trailhead) score() int {
	m := map[point]int{}
	for _, r := range th.routes {
		m[r[len(r)-1]]++
	}
	return len(m)
}

func (g *grid) findTrailheads() []*trailhead {
	var ths []*trailhead
	for y, r := range g.cells {
		for x, c := range r {
			if c == 0 {
				ths = append(ths, &trailhead{start: point{x, y}})
			}
		}
	}
	return ths
}

type dir int

const (
	up    = 0
	right = 1
	down  = 2
	left  = 3
)

func (d dir) String() string {
	switch d {
	case up:
		return "up"
	case right:
		return "right"
	case down:
		return "down"
	case left:
		return "left"
	}
	return "unknowndir"
}

type state struct {
	level   int
	todo    []dir
	visited route
}

type routeFinder struct {
	g      *grid
	states []*state
}

func newRouteFinder(g *grid) *routeFinder {
	return &routeFinder{g: g}
}

func (rf *routeFinder) addRoutesFor(th *trailhead) {
	//log.Printf("Analysing trailhead at %v height %d", th.start, rf.g.heightAt(th.start))
	// Iniialise search, retaining current and previous states in a stack.
	pos := th.start
	st := &state{
		level:   rf.g.heightAt(pos),
		visited: []point{pos},
	}
	rf.states = append(rf.states, st)

	// Start visit of a new position
	rf.iterate(pos, func(st *state) {
		th.routes = append(th.routes, st.visited)
	})
}

func (rf *routeFinder) iterate(pos point, onRouteDone func(st *state)) {
	if len(rf.states) == 0 {
		log.Fatalf("Can't iterate when state stack is empty!")
	}
	// Are we at the max height of 9? If so, report this route.
	st := rf.states[len(rf.states)-1]
	if rf.g.heightAt(pos) == 9 {
		//log.Printf("Completed route at %s height 9", pos)
		onRouteDone(st)
		return
	}

	if pos.x > 0 {
		st.todo = append(st.todo, left)
	}
	if pos.x < rf.g.w-1 {
		st.todo = append(st.todo, right)
	}
	if pos.y > 0 {
		st.todo = append(st.todo, up)
	}
	if pos.y < rf.g.h-1 {
		st.todo = append(st.todo, down)
	}
	//log.Printf("From %v can go %v", pos, st.todo)

	// Iterate todo list.
	for _, dir := range st.todo {
		var next point
		switch dir {
		case up:
			next = point{x: pos.x, y: pos.y - 1}
		case right:
			next = point{x: pos.x + 1, y: pos.y}
		case down:
			next = point{x: pos.x, y: pos.y + 1}
		case left:
			next = point{x: pos.x - 1, y: pos.y}
		}
		// Can only move to a location with height 1 greater than current height.
		if rf.g.heightAt(next) != st.level+1 {
			//log.Printf("Not going to %v because it's wrong height %d want %d", next, rf.g.heightAt(next), st.level+1)
			continue
		}
		// This height looks good. Stack new state and iterate.
		//log.Printf("Trying move from %v height %d to %v height %d", pos, st.level, next, st.level+1)
		nextSt := &state{
			level:   st.level + 1,
			visited: make([]point, len(st.visited)),
		}
		copy(nextSt.visited, st.visited)
		nextSt.visited = append(nextSt.visited, next)
		rf.states = append(rf.states, nextSt)
		rf.iterate(next, onRouteDone)
	}
}

func main() {
	log.Println("AoC-2024-day10-part1")
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
	log.Printf("Grid:\n%v", g)

	ths := g.findTrailheads()
	rf := newRouteFinder(g)
	score := 0
	for i, th := range ths {
		rf.addRoutesFor(th)
		log.Printf("Trailhead %d has %d routes (%d unique endpoints == score)", i, len(th.routes), th.score())
		score += th.score()
		/*for j, r := range th.routes {
			log.Printf(" - Route %d: %v", j, r)
		}*/
	}
	log.Printf("Overall score: %v", score)
}
