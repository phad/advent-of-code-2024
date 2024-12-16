package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

/* Example input
###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############
*/

type direction rune

const (
	north = '^'
	east  = '>'
	south = 'v'
	west  = '<'
)

func (d direction) String() string {
	return string([]rune{rune(d)})
}

type move int

const (
	unknown = iota
	cwTurn
	ccwTurn
	advance
)

func (m move) String() string {
	return map[move]string{
		unknown: "unk",
		cwTurn:  "cw",
		ccwTurn: "ccw",
		advance: "adv",
	}[m]
}

type state struct {
	// pos is the start position before the move
	pos point
	// dir is the direction the reindeer faces before the move
	dir direction
	// move is the move chosen at position pos
	mv move
	// num is how many times the move happened
	num int
}

func (s state) String() string {
	return fmt.Sprintf("[%v(%v)->%v*%d]", s.pos, s.dir, s.mv, s.num)
}

type model struct {
	arena  *grid
	states []state
	start  point
	end    point
	lowest it
}

func (m *model) String() string {
	if len(m.states) == 0 {
		fmt.Sprintf("<init>")
	}
	return fmt.Sprintf("Current:%v #prev:%d cost: %d lowest-so-far:%d", m.states[len(m.states)-1], len(m.states)-1, m.cost(), m.lowest)
}

func newModel(lines []string) (*model, error) {
	log.Printf(">>newModel")
	arena, err := newGrid(lines, true /*=wantSquare*/)
	if err != nil {
		return nil, err
	}

	startPos, ok := arena.find('S')
	if !ok {
		return nil, fmt.Errorf("Can't find start pos!")
	}

	endPos, ok := arena.find('E')
	if !ok {
		return nil, fmt.Errorf("Can't find end pos!")
	}
	log.Printf("<<newModel")
	return &model{arena: arena, start: startPos, end: endPos, lowest: 99999999999}, nil
}

// doMove() returns true if a move was made successfully.
// use atEnd() to find out if we reached the end or not.
func (m *model) doMove() bool {
	// Handle initial state
	if len(m.states) == 0 {
		m.states = append(m.states, state{
			pos: m.start,
			dir: east,
			mv:  unknown,
			num: 0,
		})
	}
	return m.innerMove()
}

// Internals of doMove, extracted for recursion.
// The state must have at least one state pushed to it.
// This will be mutated as the reindeer explores different options.
func (m *model) innerMove() bool {
	_ = readStdio()
	log.Printf("innerMove: model=%v", m)
	log.Printf("\n%v\n", m.arena)

	if m.atEnd() {
		cost := m.cost()
		if cost < m.lowest {
			m.lowest = cost
		}
		log.Printf("Reached the end at cost of: %d (lowest so far: %d)", cost, m.lowest)
		return cost < m.lowest
	}
	stateSize := len(m.states)
	if stateSize == 0 {
		log.Fatalf("invariant violated: state stack must not be empty!")
	}

	curSt := &(m.states[len(m.states)-1])

	// Figure out if a turn is available
	// Turns are very costly so prefer to advance first.
	var hasAdvance bool
	var availTurns []move
	// Fast-forward when only advances, no turns, are available.
	for {
		hasAdvance = m.canAdvance(*curSt)
		availTurns = m.canTurn(*curSt)
		if !hasAdvance || len(availTurns) > 0 {
			break
		}
		curSt.mv = advance
		curSt.num++
		nextSt := m.prepareNext(*curSt)
		curSt.pos = nextSt.pos
		_ = m.arena.set(nextSt.pos, rune(nextSt.dir))
		log.Printf("innerMove: model=%v", m)
		log.Printf("\n%v\n", m.arena)
	}
	var availMoves []move
	if hasAdvance {
		availMoves = append(availMoves, advance)
	}
	availMoves = append(availMoves, availTurns...)

	// Try each potential move. Recurse to make this DPS,
	// which will minimise cost.
	for _, mv := range availMoves {
		curSt.mv = mv
		nextSt := m.prepareNext(*curSt)
		m.states = append(m.states, nextSt)
		_ = m.arena.set(nextSt.pos, rune(nextSt.dir))
		//log.Printf("Trying move %v\nState-stack:\n%v", mv, m.states)
		if done := m.innerMove(); done {
			// The move looked ok so continue from here.
			return true
		}
		// Move 'mv' didn't work out, so unwind
		m.states = m.states[0 : len(m.states)-1]
		_ = m.arena.set(nextSt.pos, '.')
	}
	// Looks like none of the available moves worked - indicate need to backtrack.
	if len(m.states) != stateSize {
		log.Fatalf("invariant violated: state stack must be same size as previously (got %d, want %d)", len(m.states), stateSize)
	}
	return false
}

func (m *model) atEnd() bool {
	return len(m.states) > 0 && m.states[len(m.states)-1].pos == m.end
}

// Helpers

func (m *model) canAdvance(st state) bool {
	// Given st.pos and st.dir, calculate the cell we'd
	// move into.  Return false if it's a wall '#'.
	nextSt := m.prepareNext(state{pos: st.pos, dir: st.dir, mv: advance})
	nextCell, ok := m.arena.at(nextSt.pos)
	if !ok {
		log.Fatalf("Ran off the grid at %v!", nextSt.pos)
	}
	// Check if we've been to nextCell before - if so, avoid.
	for _, prev := range m.states {
		if prev.pos == nextSt.pos && prev.dir == nextSt.dir {
			log.Printf("Not visiting previously visited cell %v in same direction %v", prev.pos, prev.dir)
			return false
		}
	}
	return nextCell != '#'
}

func (m *model) canTurn(st state) []move {
	// Given st.pos and st.dir, calculate if neighbouring
	// cells (other than the one 'advance' goes to) are
	// available.
	var cwPos, ccwPos point
	switch st.dir {
	case north:
		cwPos, ccwPos = point{st.pos.x + 1, st.pos.y}, point{st.pos.x - 1, st.pos.y}
	case east:
		cwPos, ccwPos = point{st.pos.x, st.pos.y + 1}, point{st.pos.x, st.pos.y - 1}
	case south:
		cwPos, ccwPos = point{st.pos.x - 1, st.pos.y}, point{st.pos.x + 1, st.pos.y}
	case west:
		cwPos, ccwPos = point{st.pos.x, st.pos.y - 1}, point{st.pos.x, st.pos.y + 1}
	}
	var moves []move
	if cwCell, ok := m.arena.at(cwPos); ok && cwCell != '#' {
		moves = append(moves, cwTurn)
	}
	if ccwCell, ok := m.arena.at(ccwPos); ok && ccwCell != '#' {
		moves = append(moves, ccwTurn)
	}
	return moves
}

func (m *model) prepareNext(currSt state) state {
	next := state{pos: currSt.pos, dir: currSt.dir, mv: unknown}
	if currSt.mv == advance {
		switch currSt.dir {
		case north:
			next.pos.y--
		case east:
			next.pos.x++
		case south:
			next.pos.y++
		case west:
			next.pos.x--
		}
		return next
	}
	if currSt.mv == cwTurn {
		cwRotates := map[direction]direction{
			north: east,
			east:  south,
			south: west,
			west:  north,
		}
		next.dir = cwRotates[currSt.dir]
		return next
	}
	if currSt.mv == ccwTurn {
		ccwRotates := map[direction]direction{
			north: west,
			west:  south,
			south: east,
			east:  north,
		}
		next.dir = ccwRotates[currSt.dir]
		return next
	}
	// unknown - no change.
	return currSt
}

func (m *model) cost() int {
	sum := 0
	for _, st := range m.states {
		if st.mv == cwTurn || st.mv == ccwTurn {
			sum += 1000
			continue
		}
		if st.mv == advance {
			sum += st.num
		}
	}
	return sum
}

func readStdio() string {
	reader := bufio.NewReader(os.Stdin)
	in, _ := reader.ReadString('\n')
	return in
}

func main() {
	log.Println("AoC-2024-day16-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	m, err := newModel(lines)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for {
		if ok := m.doMove(); !ok {
			log.Printf("doMove()=false: we're probably not done yet?")
			break
		}
		if m.atEnd() {
			log.Printf("Detected that we reached the end!")
			break
		}
	}

	log.Printf("Cost: %d", m.cost())
}
