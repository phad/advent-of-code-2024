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
	x, y int
	dir  entity
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
		w:        9,
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

func (a *arena) String() string {
	s := ""
	for _, r := range a.entities {
		s += string(r) + "\n"
	}
	s += fmt.Sprintf("Guard: %v\n", a.g)
	return s
}

func (a *arena) step() (int, bool) {
	next := guard{x: a.g.x, y: a.g.y, dir: a.g.dir}
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

func main() {
	log.Println("AoC-2024-day06-part1")
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

	numVisited := 0
	lastState := a.String()
	for {
		log.Printf("Arena:\n%v", a)
		num, done := a.step()
		newState := a.String()
		if newState == lastState {
			log.Fatal("Stuck!!")
		}
		lastState = newState
		if done {
			numVisited = num
			break
		}
	}
	log.Printf("Guard visited %d locations", numVisited)
}
