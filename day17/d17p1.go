package main

import (
	"fmt"
	"log"
	"os"
)

/* Example input
Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0
*/

func assert(s string, b bool) {
	if !b {
		log.Fatalf("boom: %v", s)
	}
}

func runTests() {
	// If register C contains 9, the program 2,6 would set register B to 1.
	{
		c := initComputer(0, 0, 9, []operation{{2, 6}})
		err := c.execute()
		assert(fmt.Sprintf("exec err %v", err), err == nil)
		assert(fmt.Sprintf("c.B is %d want 1", c.B), c.B == 1)
	}
	// If register A contains 10, the program 5,0,5,1,5,4 would output 0,1,2.
	{
		c := initComputer(10, 0, 0, []operation{{5, 0}, {5, 1}, {5, 4}})
		err := c.execute()
		assert(fmt.Sprintf("exec err %v", err), err == nil)
		out := c.out()
		assert(fmt.Sprintf("out is %v want 0,1,2", out), out == "0,1,2")
	}
	// If register A contains 2024, the program 0,1,5,4,3,0 would output
	// 4,2,5,6,7,7,7,7,3,1,0 and leave 0 in register A
	{
		c := initComputer(2024, 0, 0, []operation{{0, 1}, {5, 4}, {3, 0}})
		err := c.execute()
		assert(fmt.Sprintf("exec err %v", err), err == nil)
		out := c.out()
		want := "4,2,5,6,7,7,7,7,3,1,0"
		assert(fmt.Sprintf("out is %v want %v", out, want), out == want)
	}
	// If register B contains 29, the program 1,7 would set register B to 26.
	{
		c := initComputer(0, 29, 0, []operation{{1, 7}})
		err := c.execute()
		assert(fmt.Sprintf("exec err %v", err), err == nil)
		assert(fmt.Sprintf("c.B is %d want 26", c.B), c.B == 26)
	}
	// If register B contains 2024 and register C contains 43690, the program
	// 4,0 would set register B to 44354
	{
		c := initComputer(0, 2024, 43690, []operation{{4, 0}})
		err := c.execute()
		assert(fmt.Sprintf("exec err %v", err), err == nil)
		assert(fmt.Sprintf("c.B is %d want 44354", c.B), c.B == 44354)
	}
	// Example
	{
		c := initComputer(729, 0, 0, []operation{{0, 1}, {5, 4}, {3, 0}})
		err := c.execute()
		assert(fmt.Sprintf("exec err %v", err), err == nil)
		out := c.out()
		want := "4,6,3,5,6,3,5,2,1,0"
		assert(fmt.Sprintf("out is %v want %v", out, want), out == want)
	}
	// Input; part 1
	// Register A: 30878003
	// Register B: 0
	// Register C: 0
	// Program: 2,4,1,2,7,5,0,3,4,7,1,7,5,5,3,0
	{
		c := initComputer(30878003, 0, 0, []operation{{2, 4}, {1, 2}, {7, 5}, {0, 3}, {4, 7}, {1, 7}, {5, 5}, {3, 0}})
		err := c.execute()
		assert(fmt.Sprintf("exec err %v", err), err == nil)
		out := c.out()
		want := "7,1,3,7,5,1,0,3,4"
		assert(fmt.Sprintf("out is %v want %v", out, want), out == want)
	}

}

func main() {
	log.Println("AoC-2024-day17-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	runTests()

	log.Printf("Input: %v", lines)
	c, err := parseInput(lines)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Computer initial state: %v", c)

	if err = c.execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Execution complete; output=%v", c.out())
}
