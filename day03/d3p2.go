package main

import (
	"log"
	"os"
	"regexp"
)

var lineRE = regexp.MustCompile("don't|do|mul\\(([0-9]+),([0-9]+)\\)")

func main() {
	log.Println("AoC-2024-day03-part2")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	total := int64(0)
	enabled := true
	for idx, line := range lines {
		matches := lineRE.FindAllStringSubmatch(line, -1)
		log.Printf("\n#%d: %q\n->%v", idx, line, matches)
		for _, m := range matches {
			log.Printf("Next match: %q", m[0])
			if m[0] == "do" {
				log.Printf("Enabling! (was enabled=%t)", enabled)
				enabled = true
				continue
			} else if m[0] == "don't" {
				log.Printf("Disabling! (was enabled=%t)", enabled)
				enabled = false
			}
			if !enabled {
				continue
			}
			a, b := mustParseInt(m[1]), mustParseInt(m[2])
			log.Printf("match: %q %dx%d=%d", m[0], a, b, a*b)
			total += a * b
		}
	}
	log.Printf("Total of all matching mul()s: %d", total)
}
