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

func main() {
	log.Println("AoC-2024-day11-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if len(lines) != 1 {
		log.Fatalf("Too many lines: %d want 1", len(lines))
	}

	bits := strings.Split(lines[0], " ")

	var seq []int
	for _, n := range bits {
		seq = append(seq, int(mustParseInt(n)))
	}

	for it := 0; it < 25; it++ {
		log.Printf("Iter %d: current seq: %v", it, seq)
		var next []int
		for _, val := range seq {
			// Rule 1: If the stone is engraved with the number 0,
			// it is replaced by a stone engraved with the number 1.
			if val == 0 {
				next = append(next, 1)
				continue
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
				continue
			}
			// Otherwise: the stone is replaced by a new stone; the
			// old stone's number multiplied by 2024 is engraved on
			// the new stone.
			next = append(next, 2024*val)
		}
		seq = next
	}
	log.Printf("Final seq:\n%v\nlength %d", seq, len(seq))
}
