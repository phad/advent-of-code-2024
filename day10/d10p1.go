package main

import (
	"log"
	"os"

	"github.com/phad/advent-of-code-2024/aoc"
)

func heightAt(g *aoc.Grid[int], p aoc.Point) int {
	v, ok := g.At(p)
	if !ok {
		log.Fatalf("Point %v is outside the grid!", p)
	}
	return v
}

type route []aoc.Point

type trailhead struct {
	start  aoc.Point
	routes []route
}

func (th trailhead) score() int {
	m := map[aoc.Point]int{}
	for _, r := range th.routes {
		m[r[len(r)-1]]++
	}
	return len(m)
}

func findTrailheads(g *aoc.Grid[int]) []*trailhead {
	var ths []*trailhead
	for y, r := range g.Cells {
		for x, c := range r {
			if c == 0 {
				ths = append(ths, &trailhead{start: aoc.Point{x, y}})
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
	g      *aoc.Grid[int]
	states []*state
}

func newRouteFinder(g *aoc.Grid[int]) *routeFinder {
	return &routeFinder{g: g}
}

func (rf *routeFinder) addRoutesFor(th *trailhead) {
	//log.Printf("Analysing trailhead at %v height %d", th.start, heightAt(rf.g, th.start))
	// Iniialise search, retaining current and previous states in a stack.
	pos := th.start
	st := &state{
		level:   heightAt(rf.g, pos),
		visited: []aoc.Point{pos},
	}
	rf.states = append(rf.states, st)

	// Start visit of a new position
	rf.iterate(pos, func(st *state) {
		th.routes = append(th.routes, st.visited)
	})
}

func (rf *routeFinder) iterate(pos aoc.Point, onRouteDone func(st *state)) {
	if len(rf.states) == 0 {
		log.Fatalf("Can't iterate when state stack is empty!")
	}
	// Are we at the max height of 9? If so, report this route.
	st := rf.states[len(rf.states)-1]
	if heightAt(rf.g, pos) == 9 {
		//log.Printf("Completed route at %s height 9", pos)
		onRouteDone(st)
		return
	}

	if pos.X > 0 {
		st.todo = append(st.todo, left)
	}
	if pos.X < rf.g.W-1 {
		st.todo = append(st.todo, right)
	}
	if pos.Y > 0 {
		st.todo = append(st.todo, up)
	}
	if pos.Y < rf.g.H-1 {
		st.todo = append(st.todo, down)
	}
	//log.Printf("From %v can go %v", pos, st.todo)

	// Iterate todo list.
	for _, dir := range st.todo {
		var next aoc.Point
		switch dir {
		case up:
			next = aoc.Point{X: pos.X, Y: pos.Y - 1}
		case right:
			next = aoc.Point{X: pos.X + 1, Y: pos.Y}
		case down:
			next = aoc.Point{X: pos.X, Y: pos.Y + 1}
		case left:
			next = aoc.Point{X: pos.X - 1, Y: pos.Y}
		}
		// Can only move to a location with height 1 greater than current height.
		if heightAt(rf.g, next) != st.level+1 {
			//log.Printf("Not going to %v because it's wrong height %d want %d", next, heightAt(rf.g, next), st.level+1)
			continue
		}
		// This height looks good. Stack new state and iterate.
		//log.Printf("Trying move from %v height %d to %v height %d", pos, st.level, next, st.level+1)
		nextSt := &state{
			level:   st.level + 1,
			visited: make([]aoc.Point, len(st.visited)),
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
	lines, err := aoc.ReadLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	g, err := aoc.NewDigitGrid(lines, true /*=wantSquare*/)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Grid:\n%v", g)

	ths := findTrailheads(g)
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
