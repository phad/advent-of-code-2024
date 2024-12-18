package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

/* Example input
Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0

but a glitch changes the initial register A value
find the new A that causes program to output itself.
*/

func convert(s string) []int {
	part1 := strings.Split(s, ":")[1]
	seq := part1[1:len(part1)]
	bits := strings.Split(seq, ",")
	var out []int
	for _, b := range bits {
		out = append(out, int(mustParseInt(b)))
	}
	return out
}

func permute(in []int) []int {
	perms := map[int]int{
		5: 0,
		4: 1,
		// 5: 2,
		7: 3,
		1: 4,
		0: 5,
		3: 6,
		2: 7,
	}

	var out []int
	for _, v := range in {
		out = append(out, perms[v])
	}
	return out
}

func seed(seq []int) int {
	s := 0
	if len(seq) > 0 {
		for idx := len(seq) - 1; idx >= 0; idx-- {
			v := seq[idx]
			s += v * 8
			s *= 8
		}
	}
	return s / 8
}

func main() {
	log.Println("AoC-2024-day17-part2")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	input := convert(lines[len(lines)-1])

	var glitchedA int
	if len(os.Args) == 3 {
		log.Printf("Overriding glitchedA: %s", os.Args[2])
		glitchedA = int(mustParseInt(os.Args[2]))
	} else {
		glitchedA = seed(permute(input))
		log.Printf("Calculating glitchedA: %d", glitchedA)
	}

	log.Printf("\nInput: %v\nGlitched A: %d\nA binary: %v", lines, glitchedA, strconv.FormatInt(int64(glitchedA), 2))

	c, err := parseInput(lines)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	c.A = glitchedA

	log.Printf("Computer initial state: %v", c)

	if err = c.execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Computer final state: %v", c)

	output := c.out()
	log.Printf("%d: Execution complete; output=%v\ninput=%v", glitchedA, output, input)

}
