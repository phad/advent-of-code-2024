package main

import (
	"log"
	"os"
	"regexp"

	"github.com/phad/advent-of-code-2024/aoc"
)

var lineRE = regexp.MustCompile("mul\\(([0-9]+),([0-9]+)\\)")

func main() {
	log.Println("AoC-2024-day03-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := aoc.ReadLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	total := int64(0)
	for idx, line := range lines {
		matches := lineRE.FindAllStringSubmatch(line, -1)
		log.Printf("\n#%d: %q\n->%v", idx, line, matches)
		for _, m := range matches {
			a, b := aoc.MustParseInt(m[1]), aoc.MustParseInt(m[2])
			log.Printf("match: %q %dx%d=%d", m[0], a, b, a*b)
			total += a * b
		}
	}
	log.Printf("Total of all matching mul()s: %d", total)
}
