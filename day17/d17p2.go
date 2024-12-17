package main

import (
	"log"
	"os"
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

func main() {
	log.Println("AoC-2024-day17-part2")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Input: %v", lines)

	iters, glitchedA := 0, 0
	for {
		c, err := parseInput(lines)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		c.A = glitchedA

		wantOutput := lines[4][strings.Index(lines[4], ":")+2 : len(lines[4])]

		//log.Printf("Computer initial state: %v", c)

		if err = c.execute(); err != nil {
			log.Fatalf("Error: %v", err)
		}
		output := c.out()
		log.Printf("%d: Execution complete; output=%v want=%v", glitchedA, output, wantOutput)

		if output == wantOutput {
			log.Printf("Quine when A glitched to %d", glitchedA)
			break
		}

		iters++
		glitchedA++
		if len(output) < len(wantOutput)-1 {
			glitchedA *= 10
		}
		if iters == 20 {
			break
		}
	}
}
