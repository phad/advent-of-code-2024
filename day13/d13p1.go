package main

import (
	"log"
	"math"
	"os"
	"regexp"
)

/* Example input
 */

type pos struct{ x, y int }
type vec struct{ dx, dy int }
type machine struct {
	a, b         vec
	costA, costB int
	prize        pos
}

func (m machine) wins(numA, numB int) bool {
	dx := numA*m.a.dx + numB*m.b.dx
	dy := numA*m.a.dy + numB*m.b.dy
	return pos{x: dx, y: dy} == m.prize
}

func (m machine) cost(numA, numB int) int {
	return numA*m.costA + numB*m.costB
}

func withinTolerance(a, b, t float64) bool {
	return math.Abs(a-b) < t
}

func (m machine) solve() (ok bool, numA, numB int) {
	ok, numA, numB = false, 0, 0
	for b := 0; b <= 100; b++ {
		a1 := float64(m.prize.x-b*m.b.dx) / float64(m.a.dx)
		a2 := float64(m.prize.y-b*m.b.dy) / float64(m.a.dy)
		if withinTolerance(a1, a2, 1e-09) && a1 >= 0.0 {
			ok = true
			numA = int(a1)
			numB = b
			return
		}
	}
	return
}

var (
	re1 = regexp.MustCompile("X\\+(\\d+), Y\\+(\\d+)")
	re2 = regexp.MustCompile("X=(\\d+), Y=(\\d+)")
)

func mustExtractTwoInt(re *regexp.Regexp, in string) (int, int) {
	matches := re.FindAllStringSubmatch(in, -1)
	for _, m := range matches {
		//log.Printf("match #%v: %v", i, m)
		return int(mustParseInt(m[1])), int(mustParseInt(m[2]))
	}
	return 0, 0
}

func parseInput(in []string) ([]machine, error) {
	var machines []machine
	m := machine{costA: 3, costB: 1}
	for idx, line := range in {
		switch idx % 4 {
		case 0:
			x, y := mustExtractTwoInt(re1, line)
			m.a = vec{x, y}
		case 1:
			x, y := mustExtractTwoInt(re1, line)
			m.b = vec{x, y}
		case 2:
			x, y := mustExtractTwoInt(re2, line)
			m.prize = pos{x, y}
		case 3:
			log.Printf("Built machine %v", m)
			machines = append(machines, m)
		}
	}
	log.Printf("Built final machine %v", m)
	machines = append(machines, m)
	return machines, nil
}

func main() {
	log.Println("AoC-2024-day13-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	machines, err := parseInput(lines)
	if err != nil {
		log.Fatalf("Input error: %v", err)
	}

	tokens, numWon := 0, 0
	for idx, m := range machines {
		ok, numA, numB := m.solve()
		if ok {
			c := m.cost(numA, numB)
			tokens += c
			numWon++
			log.Printf("Machine #%d: won with %d A and %d B presses", idx, numA, numB)
			continue
		}
		log.Printf("Machine #%d can't be won.", idx)
	}
	log.Printf("%d of %d machines can be won for %d tokens", numWon, len(machines), tokens)
}
