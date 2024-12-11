package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

/* Example input
0 1 10 99 999
*/

func apply(val int) []int {
	var next []int
	// Rule 1: If the stone is engraved with the number 0,
	// it is replaced by a stone engraved with the number 1.
	if val == 0 {
		next = append(next, 1)
		return next
	}
	// Rule 2: If the stone is engraved with a number that
	// has an even number of digits, it is replaced by two
	// stones. The left half of the digits are engraved on
	// the new left stone, and the right half of the digits
	// are engraved on the new right stone. (The new numbers
	// don't keep extra leading zeroes: 1000 would become
	// stones 10 and 0.)
	s := fmt.Sprintf("%d", val)
	if len(s)%2 == 0 {
		next = append(next, int(mustParseInt(s[0:len(s)/2])))
		next = append(next, int(mustParseInt(s[len(s)/2:len(s)])))
		return next
	}
	// Otherwise: the stone is replaced by a new stone; the
	// old stone's number multiplied by 2024 is engraved on
	// the new stone.
	return append(next, 2024*val)
}

type production struct {
	value    int
	occurs   int
	produces []int
}

func (p *production) String() string {
	return fmt.Sprintf("%d(%d #)->%v", p.value, p.occurs, p.produces)

}

type productionSet struct {
	pr map[int]*production
}

func newProductionSet(vals []int) productionSet {
	ps := productionSet{
		pr: make(map[int]*production),
	}
	for _, v := range vals {
		ps.insert(v, 1)
	}
	return ps
}

func (ps productionSet) String() string {
	var s strings.Builder
	var vals []int
	for v, p := range ps.pr {
		if p.occurs > 0 {
			vals = append(vals, v)
		}
	}
	for _, v := range vals {
		s.WriteString(ps.pr[v].String() + "\n")
	}
	return s.String()
}

func (ps *productionSet) insert(v, n int) {
	p, ok := ps.pr[v]
	if !ok {
		p = &production{
			value:    v,
			occurs:   0,
			produces: apply(v),
		}
		ps.pr[v] = p
	}
	p.occurs += n
}

// replace all occurrences of the (unique) vals
func (ps *productionSet) replace(vals []int) {
	toInsert := map[int]int{}
	for _, v := range vals {
		p, ok := ps.pr[v]
		if !ok {
			log.Fatalf("No production for %d", v)
		}
		if p.occurs == 0 {
			log.Fatalf("Can't replace %d as it occurs 0 times", v)
		}
		for _, i := range p.produces {
			toInsert[i] += p.occurs
		}
		p.occurs = 0
	}

	for v, occ := range toInsert {
		ps.insert(v, occ)
	}
}

func (ps *productionSet) count() int {
	total := 0
	for _, p := range ps.pr {
		total += p.occurs
	}
	return total
}

func main() {
	log.Println("AoC-2024-day11-part2")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file> [<iters>]")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if len(lines) != 1 {
		log.Fatalf("Too many lines: %d want 1", len(lines))
	}

	iters := 25
	if len(os.Args) == 3 {
		iters = int(mustParseInt(os.Args[2]))
	}

	bits := strings.Split(lines[0], " ")

	var seq []int
	for _, n := range bits {
		seq = append(seq, int(mustParseInt(n)))
	}

	ps := newProductionSet(seq)

	// Iteration
	for it := 0; it < iters; it++ {
		//log.Printf("Iter:%d, have:\n%v\n", it, ps)
		log.Printf("Iter: %d", it)
		var vals []int
		for v, p := range ps.pr {
			if p.occurs > 0 {
				vals = append(vals, v)
			}
		}
		ps.replace(vals)
	}

	log.Printf("Final count: %d", ps.count())
	//log.Printf("Final:\n%v\nCount: %d", ps, ps.count())
}
