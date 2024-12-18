package main

import (
	"log"
	"os"
)

/* Example input
 */

func main() {
	log.Println("AoC-2024-day18-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Input: %v", lines)
}
