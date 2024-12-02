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
	return l.doIsSafe(1)
}

func (l level) doIsSafe(numViolsAllowed int) bool {
	log.Printf("doIsSafe%v(viols left %d)", l, numViolsAllowed)
	if len(l) < 2 {
		return false
	}
	violIdx := -1
	increasing, decreasing := l[0] < l[len(l)-1], l[0] > l[len(l)-1]
	for i := 1; i < len(l); i++ {
		change := l[i] - l[i-1]
		if increasing {
			if change <= 0 {
				violIdx = i
				break
			}
			if change > 3 {
				violIdx = i
				break
			}
		}
		if decreasing {
			if change >= 0 {
				violIdx = i
				break
			}
			if -change > 3 {
				violIdx = i
				break
			}
		}
	}
	if violIdx == -1 {
		// No violation found
		return true
	}
	log.Printf("**Found violating element #%d = %d", violIdx, l[violIdx])
	if numViolsAllowed == 0 {
		log.Printf("Too bad, no more violations allowed")
		return false
	}

	// Tolerate 1 violation. Try amending at each position.
	for idx := 0; idx < len(l); idx++ {
		log.Printf("**Try removing #%d element %d", idx, l[idx])
		l2 := make(level, 0, len(l)-1)
		l2 = append(l2, l[0:idx]...)
		l2 = append(l2, l[idx+1:]...)
		if l2.doIsSafe(numViolsAllowed - 1) {
			log.Printf("yay, dampener has fixed it")
			return true
		}
	}
	log.Printf("oh noes, no dampening possibility exists")
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
		log.Printf("level %v: safe? %t\n", l, safe)
	}
	log.Printf("#safe levels: %d", numSafe)
}
