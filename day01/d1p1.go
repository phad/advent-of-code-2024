package main

import (
	"log"
	"math"
	"regexp"
	"sort"

	"github.com/phad/advent-of-code-2024/aoc"
)

var lineRE = regexp.MustCompile("([0-9]+)")

func main() {
	lines := aoc.MustReadInput("AoC-2024-day01-part1")

	// Two slices of integers read from the input file.
	var left, right []int64

	// Each line is formatted as `<number><whitespace><number>`
	for idx, line := range lines {
		matches := lineRE.FindAllString(line, -1)
		if len(matches) != 2 {
			log.Fatalf("Error: input line %d %q did not contain two numbers.", idx, line)
		}
		// log.Printf("Input line %d contains %v", idx, matches)
		left = append(left, aoc.MustParseInt(matches[0]))
		right = append(right, aoc.MustParseInt(matches[1]))
	}
	// log.Printf("Left: %v", left)
	// log.Printf("Right %v", right)

	// Sort left and right, then we can measure distances
	sort.Slice(left, func(i, j int) bool { return left[i] < left[j] })
	sort.Slice(right, func(i, j int) bool { return right[i] < right[j] })
	dist := int64(0)
	for idx, l := range left {
		r := right[idx]
		d := int64(math.Abs(float64(r - l)))
		// log.Printf("l: %d r: %d: d: %d", l, r, d)
		dist += d
	}
	log.Printf("Overall distance: %d", dist)

}
