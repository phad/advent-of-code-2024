// Package aoc provides small helpers shared across the Advent of Code 2024
// solutions. It collects logic that was previously copy-pasted into a per-day
// util.go in every day directory.
package aoc

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

// ReadLines opens the file at path f and returns its contents split into lines.
func ReadLines(f string) ([]string, error) {
	rd, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer rd.Close()

	var lines []string
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// MustParseInt parses s as a base-10 int64, exiting the program via log.Fatalf
// if the value cannot be parsed.
func MustParseInt(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("ParseInt(%q) err=%v", s, err)
	}
	return v
}

// MustReadInput reads the input file named by the first command-line argument,
// logging a banner with the given label. It exits the program via log.Fatal if
// no argument is supplied or the file cannot be read. It centralises the
// boilerplate previously repeated at the top of every day's main().
func MustReadInput(label string) []string {
	log.Println(label)
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := ReadLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Read %d input lines", len(lines))
	return lines
}
