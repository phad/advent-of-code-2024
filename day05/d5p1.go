package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

/*
* Example input: one or more in each group; groups separated by empty line.
53|13

75,47,61,53,29
*/

type ruleSet struct {
	orderingRules    map[int]map[int]bool
	revOrderingRules map[int]map[int]bool
}

func newRuleSet() *ruleSet {
	return &ruleSet{
		orderingRules:    map[int]map[int]bool{},
		revOrderingRules: map[int]map[int]bool{},
	}
}

func (rs *ruleSet) addOrdering(first, second int) {
	if _, ok := rs.orderingRules[first]; !ok {
		rs.orderingRules[first] = make(map[int]bool)
	}
	rs.orderingRules[first][second] = true

	if _, ok := rs.revOrderingRules[second]; !ok {
		rs.revOrderingRules[second] = make(map[int]bool)
	}
	rs.revOrderingRules[second][first] = true
}

func (rs *ruleSet) String() string {
	return fmt.Sprintf("forward ordering rules:\n%vreverse ordering rules:\n%v", rs.orderingRules, rs.revOrderingRules)
}

func (rs *ruleSet) filter(updates [][]int) []int {
	var validIndices []int
	for idx, update := range updates {
		if rs.isValid(update) {
			validIndices = append(validIndices, idx)
		}
	}
	return validIndices
}

func (rs *ruleSet) isValid(update []int) bool {
	// TODO: decide if valid
	log.Printf("Considering %v", update)
	for i, elem := range update {
		for j := i + 1; j < len(update); j++ {
			if !rs.orderingRules[elem][update[j]] {
				return false
			}
		}
	}
	return true
}

func main() {
	log.Println("AoC-2024-day05-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	rs := newRuleSet()
	updates := [][]int{}

	for _, line := range lines {
		// Try parsing as an ordering rule first
		bits := strings.Split(line, "|")
		if len(bits) == 2 {
			first := int(mustParseInt(bits[0]))
			second := int(mustParseInt(bits[1]))
			rs.addOrdering(first, second)
			continue
		}
		bits = strings.Split(line, ",")
		if len(bits) >= 2 {
			var update []int
			for _, bit := range bits {
				update = append(update, int(mustParseInt(bit)))
			}
			log.Printf("Read update sequence: %v", update)
			updates = append(updates, update)
		}

	}

	log.Printf("rules:\n%v", rs)
	log.Printf("Proposed updates:\n%v", updates)

	validUpdateIndices := rs.filter(updates)
	sum := 0
	for _, validIdx := range validUpdateIndices {
		valid := updates[validIdx]
		middlePage := valid[(len(valid)-1)/2]
		log.Printf("Valid update %v: taking middle page number %d", valid, middlePage)
		sum += middlePage
	}

	log.Printf("Sum of middle page numbers for valid updates: %d", sum)
}
