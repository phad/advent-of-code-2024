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

func (l level) hasViolation() bool {
	increasing, decreasing := l[0] < l[1], l[0] > l[1]
	if !increasing && !decreasing {
		return true
	}
	for i := 1; i < len(l); i++ {
		change := l[i] - l[i-1]
		if change == 0 {
			return true
		}
		if change > 0 && decreasing {
			return true
		}
		if change < 0 && increasing {
			return true
		}
		if change > 3 || change < -3 {
			return true
		}
	}
	return false
}

func (l level) isSafe() bool {
	log.Printf("doIsSafe%v", l)
	if len(l) < 2 {
		return false
	}
	if !l.hasViolation() {
		log.Printf("  No violations :D\n")
		return true
	}
	// Tolerate 1 violation. Try amending at each position.
	for idx := 0; idx < len(l); idx++ {
		log.Printf(" **Try removing #%d element %d", idx, l[idx])
		l2 := make(level, 0, len(l)-1)
		l2 = append(l2, l[0:idx]...)
		l2 = append(l2, l[idx+1:]...)
		if !l2.hasViolation() {
			log.Printf("  yay, dampener has fixed it\n")
			return true
		}
	}
	log.Printf(" oh noes, no dampening possibility exists\n")
	return false
}

func main() {
	log.Println("AoC-2024-day02-part2")
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
		log.Printf("level %v: safe? %t\n\n\n", l, safe)
	}
	log.Printf("#safe levels: %d", numSafe)
}
