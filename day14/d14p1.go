package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

/* Example input
p=18,3 v=-20,-92
*/

type pos struct{ x, y int }
type vec struct{ dx, dy int }

type robot struct {
	p pos
	v vec
}

func (r *robot) String() string {
	return fmt.Sprintf("p<%v> v<%v>", r.p, r.v)
}

func (r *robot) move(a arena) {
	x := r.p.x + r.v.dx
	for x < 0 || x >= a.w {
		if x < 0 {
			x += a.w
		} else if x >= a.w {
			x -= a.w
		}
	}
	y := r.p.y + r.v.dy
	for y < 0 || y >= a.h {
		if y < 0 {
			y += a.h
		} else if y >= a.h {
			y -= a.h
		}
	}
	r.p = pos{x: x, y: y}
}

type arena struct {
	w, h int
}

var (
	re = regexp.MustCompile("p=(\\d+),(\\d+) v=([-\\d]+),([-\\d]+)")
)

func mustExtractFourInts(re *regexp.Regexp, in string) (pos, vec) {
	matches := re.FindAllStringSubmatch(in, -1)
	for _, m := range matches {
		//log.Printf("match #%v: %v", i, m)
		p := pos{int(mustParseInt(m[1])), int(mustParseInt(m[2]))}
		v := vec{int(mustParseInt(m[3])), int(mustParseInt(m[4]))}
		return p, v
	}
	return pos{0, 0}, vec{0, 0}
}

func parseInput(in []string) []*robot {
	var robots []*robot
	for _, line := range in {
		p, v := mustExtractFourInts(re, line)
		robots = append(robots, &robot{p: p, v: v})
	}
	return robots
}

func debugString(iter int, robots []*robot, a arena) string {
	var s strings.Builder
	if iter == 0 {
		s.WriteString("Initial state:\n")
	} else {
		s.WriteString(fmt.Sprintf("After %d second(s)\n", iter))
	}
	rs := map[int][]*robot{}
	for _, r := range robots {
		rs[r.p.y] = append(rs[r.p.y], r)
	}
	for y := 0; y < a.h; y++ {
		sort.Slice(rs[y], func(i, j int) bool { return rs[y][i].p.x < rs[y][j].p.x })
		var rns []rune
		for x := 0; x < a.w; x++ {
			numRobots := 0
			for _, r := range rs[y] {
				if r.p.x == x {
					numRobots++
				}
			}
			if numRobots > 0 {
				rns = append(rns, rune(48+numRobots))
			} else {
				rns = append(rns, '.')
			}
		}
		s.WriteString(fmt.Sprintf("%s\n", string(rns)))
	}
	return s.String()
}

func safetyFactor(robots []*robot, a arena) int {
	c := map[bool]map[bool]int{false: map[bool]int{}, true: map[bool]int{}}
	for _, r := range robots {
		if r.p.x == a.w/2 || r.p.y == a.h/2 {
			log.Printf("Robot: %v on the boundary - skipping.", r)
			continue
		}
		isLeft := r.p.x < a.w/2
		isTop := r.p.y < a.h/2
		log.Printf("Robot: %v isLeft: %t isTop: %t", r, isLeft, isTop)
		c[isLeft][isTop]++
	}
	log.Printf("quadrant counts: %v", c)
	f := 1
	f *= c[true][true]
	f *= c[false][true]
	f *= c[false][false]
	f *= c[true][false]
	return f
}

func main() {
	log.Println("AoC-2024-day14-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	robots := parseInput(lines)
	if err != nil {
		log.Fatalf("Input error: %v", err)
	}

	log.Printf("Read %d robots", len(lines))

	var a arena
	if len(lines) < 100 {
		a = arena{w: 11, h: 7}
	} else {
		a = arena{w: 101, h: 103}
	}

	for tick := 0; tick < 100; tick++ {
		log.Printf("\n%s\n%v\n", debugString(tick, robots, a), "") //robots)
		for i := range robots {
			robots[i].move(a)
		}
	}
	log.Printf("\n%s\n%v\n", debugString(100, robots, a), "") //robots)

	log.Printf("After simulation, robots are:\n%v", robots)
	log.Printf("Safety factor: %d", safetyFactor(robots, a))
}
