package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/phad/advent-of-code-2024/aoc"
)

/* Example input
########
#..O.O.#
##@.O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

<^^>>>vv<v>>v<<
*/

type move rune

const (
	up    = '^'
	right = '>'
	down  = 'v'
	left  = '<'
)

func (m move) String() string {
	return string([]rune{rune(m)})
}

func readMoves(in []string) ([]move, error) {
	var moves []move
	for _, l := range in {
		for i := 0; i < len(l); i++ {
			m := string([]rune{rune(l[i])})
			if !strings.Contains("^>v<", m) {
				return nil, fmt.Errorf("Invalid move %v", m)
			}
			moves = append(moves, move(l[i]))
		}
	}
	return moves, nil
}

type model struct {
	arena *aoc.Grid[rune]
	moves []move
	next  int
	pos   aoc.Point
}

func (m *model) String() string {
	return fmt.Sprintf("%v\nRobot at %v\n%d moves: done %d, todo %d", m.arena, m.pos, len(m.moves), m.next, len(m.moves)-m.next)
}

func newModel(lines []string) (*model, error) {
	// Split input into grid and moves.
	dividerPos := -1
	for i, l := range lines {
		if len(l) == 0 {
			dividerPos = i
			break
		}
	}
	if dividerPos < 0 {
		return nil, fmt.Errorf("Didn't find the empty divider line.")
	}

	arena, err := aoc.NewRuneGrid(lines[0:dividerPos], true /*=wantSquare*/)
	if err != nil {
		return nil, err
	}

	moves, err := readMoves(lines[dividerPos+1 : len(lines)])
	if err != nil {
		return nil, err
	}

	pos, ok := arena.Find('@')
	if !ok {
		return nil, fmt.Errorf("Can't find robot!")
	}

	return &model{arena: arena, moves: moves, pos: pos}, nil
}

// returns true if more moves are available.
func (m *model) doMove() bool {
	if m.next == len(m.moves) {
		return false
	}
	move := m.moves[m.next]
	m.next++
	log.Printf("Moving: %v", move)

	nextPos, ok := m.innerMove(m.pos, move)
	if ok {
		m.pos = nextPos
	} else {
		//log.Printf("Couldn't move")
	}
	return true
}

// returns true if something was moved.
func (m *model) innerMove(pos aoc.Point, move move) (aoc.Point, bool) {
	var nextCell rune
	var nextPos aoc.Point
	switch move {
	case up:
		nextPos = aoc.Point{pos.X, pos.Y - 1}
	case right:
		nextPos = aoc.Point{pos.X + 1, pos.Y}
	case down:
		nextPos = aoc.Point{pos.X, pos.Y + 1}
	case left:
		nextPos = aoc.Point{pos.X - 1, pos.Y}
	}
	nextCell, ok := m.arena.At(nextPos)
	if !ok {
		log.Fatalf("Ran off the grid at %v!", nextPos)
	}
	if nextCell == '#' {
		// boundary or obstacle, can't move here.
		//log.Printf("Hit boundary trying to move to %v currently occupied by %v", nextPos, nextCell)
		return aoc.Point{}, false
	}
	if nextCell == 'O' {
		// Need to see if we can shift this first.
		if _, ok := m.innerMove(nextPos, move); !ok {
			return aoc.Point{}, false
		}
	}
	// Make the move!
	//log.Printf("Trying to swap grid cells %v<->%v!", pos, nextPos)
	if ok := m.arena.Swap(pos, nextPos); !ok {
		log.Fatalf("Failed to swap grid cells %v<->%v!", pos, nextPos)
	}
	//log.Printf("innerMove: %v", m.arena)
	return nextPos, true
}

func (m *model) gpsSum() int {
	sum := 0
	m.arena.FindAll('O', func(p aoc.Point) bool {
		coord := 100*p.Y + p.X
		sum += coord
		return true // keep going
	})
	return sum
}

func main() {
	log.Println("AoC-2024-day15-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := aoc.ReadLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	m, err := newModel(lines)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for {
		log.Printf("%v", m)
		if ok := m.doMove(); !ok {
			break
		}
	}

	log.Printf("GPS Coords Sum: %d", m.gpsSum())
}
