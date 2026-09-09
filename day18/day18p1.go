package main

import (
	"log"
	"os"

	"github.com/phad/advent-of-code-2024/aoc"
)

/* Example input
 */

func main() {
	log.Println("AoC-2024-day18-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := aoc.ReadLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Input: %v", lines)
}
