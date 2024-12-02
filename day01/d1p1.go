package main

import (
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	_ "strconv"
)

var lineRE = regexp.MustCompile("([0-9]+)")

func mustParseInt(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("ParseInt(%q) err=%v", s, err)
	}
	return v
}

func main() {
	log.Println("AoC-2024-day01-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Read %d input lines", len(lines))

	// Two slices of integers read from the input file.
	var left, right []int64

	// Each line is formatted as `<number><whitespace><number>`
	for idx, line := range lines {
		matches := lineRE.FindAllString(line, -1)
		if len(matches) != 2 {
			log.Fatalf("Error: input line %d %q did not contain two numbers.", idx, line)
		}
		// log.Printf("Input line %d contains %v", idx, matches)
		left = append(left, mustParseInt(matches[0]))
		right = append(right, mustParseInt(matches[1]))
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
