package main

import (
	"log"
	"math"
	"os"
	"regexp"
)

/* Example input
Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=10000000008400, Y=10000000005400

Button A: X+26, Y+66
etc.
*/

const (
	tol        = float64(1e-3)
	adjustment = int64(10000000000000)
)

type pos struct{ x, y int64 }
type vec struct{ dx, dy int64 }
type machine struct {
	// Vector claw moves through when buttons A and B pressed.
	a, b vec
	// Cost in tokens for pressingn buttons A and B
	costA, costB int64
	// Location of prize
	p pos
}

func (m machine) wins(numA, numB int64) bool {
	dx := numA*m.a.dx + numB*m.b.dx
	dy := numA*m.a.dy + numB*m.b.dy
	return pos{x: dx, y: dy} == m.p
}

func (m machine) cost(numA, numB int64) int64 {
	return numA*m.costA + numB*m.costB
}

func withinTolerance(a, b, t float64) bool {
	return math.Abs(a-b) < t
}

func (m machine) solve() (ok bool, numA, numB int64) {
	ok, numA, numB = false, 0, 0

	// Assume we're going to make approximately
	// adjustment/100 presses of both A, and B,
	// since no input has dx or dy > 100 for A or B.
	// We can then offset the whole problem.
	numOffset := adjustment / 100
	prizeOffset := pos{m.p.x - adjustment, m.p.y - adjustment}

}

func (m machine) solveOld() (ok bool, numA, numB int) {
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

// This method with floats should work but I think doesn't due to loss
// of precision with 64-bit floats.
func (m machine) solveFloat() (ok bool, numA, numB int64) {
	ok, numA, numB = false, 0, 0
	/*
		(m.p.x - b*m.b.dx)/m.a.dx == (m.p.y - b*m.b.dy)/m.a.dy

		m.p.x/m.a.dx - b * m.b.dx/m.a.dx == m.p.y/m.a.dy - b * m.b.dy/m.a.dy

		b * (m.b.dy/m.a.dy - m.b.dx/m.a.dx) == m.p.y/m.a.dy - m.p.x/m.a.dx

		b == (m.p.y/m.a.dy - m.p.x/m.a.dx) / (m.b.dy/m.a.dy - m.b.dx/m.a.dx)
	*/
	mpx, mpy := float64(m.p.x), float64(m.p.y)
	madx, mady := float64(m.a.dx), float64(m.a.dy)
	mbdx, mbdy := float64(m.b.dx), float64(m.b.dy)

	var b float64 = (mpy/mady - mpx/madx) / (mbdy/mady - mbdx/madx)
	var a1 float64 = (mpx - b*mbdx) / madx
	var a2 float64 = (mpy - b*mbdy) / mady

	log.Printf("\na1=%v\na2=%v\n b=%v\n", a1, a2, b)

	if withinTolerance(a1, a2, tol) && withinTolerance(float64(int64(a1)), a1, tol) && withinTolerance(float64(int64(b)), b, tol) {
		numA = int64(a1)
		numB = int64(b)
		ok = true
	}
	return
}

var (
	re1 = regexp.MustCompile("X\\+(\\d+), Y\\+(\\d+)")
	re2 = regexp.MustCompile("X=(\\d+), Y=(\\d+)")
)

func mustExtractTwoInt(re *regexp.Regexp, in string) (int64, int64) {
	matches := re.FindAllStringSubmatch(in, -1)
	for _, m := range matches {
		//log.Printf("match #%v: %v", i, m)
		return mustParseInt(m[1]), mustParseInt(m[2])
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
			m.p = pos{x + adjustment, y + adjustment}
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
	log.Println("AoC-2024-day13-part2")
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

	tokens, numWon := int64(0), 0
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
