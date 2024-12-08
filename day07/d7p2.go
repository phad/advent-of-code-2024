package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

/* input format
190: 10 19
*/

type op rune

const (
	unknown = '?'
	add     = '+'
	mult    = '*'
	concat  = '|'
)

type calc struct {
	total int64
	vals  []int64
	ops   []op
	valid bool
}

func newCalc(in string) (*calc, error) {
	bits := strings.Split(in, ":")
	if len(bits) != 2 {
		return nil, fmt.Errorf("malformed input: want <v>:<v>+, got %s", in)
	}
	vals := strings.Split(strings.TrimSpace(bits[1]), " ")
	if len(vals) < 2 {
		return nil, fmt.Errorf("malformed input: want <v>:<v>+, want >=2 vals got %d", len(vals))
	}
	c := &calc{total: mustParseInt(bits[0])}
	for i, v := range vals {
		c.vals = append(c.vals, mustParseInt(v))
		if i > 0 {
			c.ops = append(c.ops, unknown)
		}
	}
	return c, nil
}

func (c *calc) validOps() bool {
	for i := 0; i < int(math.Pow(3, float64(len(c.ops)))); i++ {
		var try []op
		j := i
		for k := 0; k < len(c.ops); k++ {
			switch j % 3 {
			case 0:
				try = append([]op{add}, try...)
			case 1:
				try = append([]op{mult}, try...)
			case 2:
				try = append([]op{concat}, try...)
			}
			j -= (j % 3)
			j /= 3
		}
		//log.Printf("Trying: %v", string(try))
		tot := c.vals[0]
		for i, o := range try {
			v := c.vals[i+1]
			switch o {
			case add:
				tot += v
			case mult:
				tot *= v
			case concat:
				tot = mustParseInt(fmt.Sprintf("%d%d", tot, v))
			}
		}
		if tot == c.total {
			c.ops = try
			c.valid = true
			break
		}
	}
	return c.valid
}

func (c *calc) String() string {
	s := fmt.Sprintf("%d:", c.total)
	for i, v := range c.vals {
		s += fmt.Sprintf(" %d", v)
		if i < len(c.vals)-1 {
			s += fmt.Sprintf(" %v", string(c.ops[i]))
		}
	}
	return s
}

func main() {
	log.Println("AoC-2024-day07-part2")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	var calcs []*calc

	for idx, line := range lines {
		c, err := newCalc(line)
		if err != nil {
			log.Fatalf("Line %d: malformed input %q", idx, line)
		}
		calcs = append(calcs, c)
	}

	total := int64(0)
	for idx, c := range calcs {
		// log.Printf("#%d: checking %v", idx, c)
		if c.validOps() {
			log.Printf("Calc %d: Valid ops: %v", idx, c)
			total += c.total
		}
	}

	log.Printf("Total for valid calculations: %d", total)
}
