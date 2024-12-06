package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type entity rune

const (
	empty      entity = '.'
	obstacle   entity = '#'
	guardUp    entity = '^'
	guardRight entity = '>'
	guardDown  entity = 'v'
	guardLeft  entity = '<'
	visited    entity = 'X'

	validEntities = ".#^>v<"
)

var rotations = map[entity]entity{
	guardUp:    guardRight,
	guardRight: guardDown,
	guardDown:  guardLeft,
	guardLeft:  guardUp,
}

type guard struct {
	x, y  int
	dir   entity
	moves int
}

func (g guard) String() string {
	return fmt.Sprintf("%v@(%d,%d)", string(rune(g.dir)), g.x, g.y)
}

type arena struct {
	w, h     int
	entities [][]entity
	g        guard
}

func initArena(in []string) (*arena, error) {
	a := &arena{
		w:        0,
		h:        len(in),
		entities: make([][]entity, len(in)),
	}
	for i, line := range in {
		if i == 0 {
			a.w = len(line)
		} else if len(line) != a.w {
			return nil, fmt.Errorf("Inconsistent width on row %d: got %d want %d", i, len(line), a.w)
		}
		row := make([]entity, a.w)
		for j, r := range line {
			v := strings.IndexRune(validEntities, r)
			if v == -1 {
				return nil, fmt.Errorf("not a valid entity %v, want [%s]", r, validEntities)
			}
			row[j] = entity(r)
			if v < 2 {
				// not a guard
				continue
			}
			a.g = guard{x: j, y: i, dir: entity(r)}
		}
		a.entities[i] = row
	}
	return a, nil
}

func (a *arena) asInput() []string {
	var in []string
	for _, r := range a.entities {
		in = append(in, string(r))
	}
	return in
}

func (a *arena) String() string {
	s := strings.Join(a.asInput(), "\n")
	s += fmt.Sprintf("\nGuard: %v\n", a.g)
	return s
}

func (a *arena) step() (int, bool) {
	next := guard{x: a.g.x, y: a.g.y, dir: a.g.dir, moves: a.g.moves + 1}
	a.entities[a.g.y][a.g.x] = visited

	switch a.g.dir {
	case guardUp:
		next.y = a.g.y - 1
	case guardRight:
		next.x = a.g.x + 1
	case guardDown:
		next.y = a.g.y + 1
	case guardLeft:
		next.x = a.g.x - 1
	}

	exited := next.x < 0 || next.x >= a.w || next.y < 0 || next.y >= a.h
	if !exited {
		if a.entities[next.y][next.x] == obstacle {
			a.g.dir = rotations[next.dir]
		} else {
			a.g = next
		}
		a.entities[a.g.y][a.g.x] = next.dir
	}

	numVisited := 0
	for i := 0; i < a.w; i++ {
		for j := 0; j < a.h; j++ {
			if a.entities[j][i] == visited {
				numVisited++
			}
		}
	}
	return numVisited, exited
}

func (a *arena) run() (int, int, bool) {
	numVisited := 0
	lastState := a.String()
	looped := false
	numVisitedUnchangedTimes := 0
	for {
		//log.Printf("Arena:\n%v", lastState)
		num, done := a.step()
		newState := a.String()
		if newState == lastState {
			log.Fatal("Stuck!!")
		}
		lastState = newState

		if numVisited == num {
			numVisitedUnchangedTimes++
			if numVisitedUnchangedTimes > a.w*a.h {
				looped = true
				break
			}
		} else {
			numVisitedUnchangedTimes = 0
		}

		numVisited = num

		if done {
			break
		}
	}
	return numVisited, a.g.moves, looped
}

func (a *arena) tryCreateLoop(x, y int) bool {
	if x < 0 || x >= a.w || y < 0 || y >= a.h || a.entities[y][x] != empty {
		return false
	}
	state := a.asInput()
	a.entities[y][x] = obstacle
	_, _, looped := a.run()
	a2, _ := initArena(state)
	a.entities, a.g = a2.entities, a2.g
	return looped
}

func main() {
	log.Println("AoC-2024-day06-part2")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	a, err := initArena(lines)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	tries, numLoops := 0, 0
	for j := 0; j < a.h; j++ {
		for i := 0; i < a.w; i++ {
			if (tries % 100) == 0 {
				log.Printf("%d tries %d loops found", tries, numLoops)
			}
			tries++
			if a.tryCreateLoop(i, j) {
				numLoops++
			}
		}
	}
	log.Printf("Done: %d loops found", numLoops)
}
