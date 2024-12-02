package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var lineRE = regexp.MustCompile("([0-9]+)")

func mustParseInt(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("ParseInt(%q) err=%v", s, err)
	}
	return v
}

type level []int64

func makeLevel(s string) (level, error) {
	matches := lineRE.FindAllString(s, -1)
	if len(matches) < 2 {
		return nil, fmt.Errorf("Error: input %q must contain at least two numbers.", s)
	}
	// log.Printf("Input line %d contains %v", idx, matches)A
	l := make(level, 0, len(matches))
	for _, m := range matches {
		l = append(l, mustParseInt(m))
	}
	return l, nil
}

func (l level) isSafe() bool {
	if len(l) < 2 {
		return false
	}
	if l[0] == l[1] {
		return false
	}
	increasing := l[0] < l[1]
	for i := 1; i < len(l); i++ {
		change := l[i] - l[i-1]
		if increasing {
			if change <= 0 {
				return false
			}
			if change > 3 {
				return false
			}
		} else {
			// decreasing
			if change >= 0 {
				return false
			}
			if -change > 3 {
				return false
			}
		}
	}
	return true
}

func main() {
	log.Println("AoC-2024-day02-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Read %d input lines", len(lines))

	numSafe := 0
	for idx, line := range lines {
		l, err := makeLevel(line)
		if err != nil {
			log.Fatalf("line #%d %q: err = %v", idx, line, err)
		}
		safe := l.isSafe()
		if safe {
			numSafe++
		}
		//log.Printf("level %v: safe? %t", l, safe)
	}
	log.Printf("#safe levels: %d", numSafe)
}
