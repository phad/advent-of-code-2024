package main

import (
	"fmt"
	"log"
	"os"
	"sort"
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

func (rs *ruleSet) splitValidInvalid(updates [][]int) ([]int, []int) {
	var validIndices, invalidIndices []int
	for idx, update := range updates {
		if rs.isValid(update) {
			validIndices = append(validIndices, idx)
		} else {
			invalidIndices = append(invalidIndices, idx)
		}
	}
	return validIndices, invalidIndices
}

func (rs *ruleSet) isValid(update []int) bool {
	valid := true
outer:
	for i, elem := range update {
		for j := i + 1; j < len(update); j++ {
			if !rs.orderingRules[elem][update[j]] {
				valid = false
				break outer
			}
		}
	}
	log.Printf("Considered %v valid? %t", update, valid)
	return valid
}

func (rs *ruleSet) sort(update []int) []int {
	if rs.isValid(update) {
		return update
	}
	sorted := make([]int, len(update))
	copy(sorted, update)
	sort.Slice(sorted, func(i, j int) bool {
		return rs.orderingRules[sorted[i]][sorted[j]]
	})
	return sorted
}

func main() {
	log.Println("AoC-2024-day05-part2")
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

	validUpdateIndices, invalidUpdateIndices := rs.splitValidInvalid(updates)
	sumValid, sumInvalid := 0, 0
	for _, validIdx := range validUpdateIndices {
		valid := updates[validIdx]
		middlePage := valid[(len(valid)-1)/2]
		log.Printf("Valid update %v: taking middle page number %d", valid, middlePage)
		sumValid += middlePage
	}
	for _, invalidIdx := range invalidUpdateIndices {
		invalid := updates[invalidIdx]
		fixed := rs.sort(invalid)
		middlePage := fixed[(len(fixed)-1)/2]
		log.Printf("Invalid update %v, fixed as %v: taking middle page number %d", invalid, fixed, middlePage)
		sumInvalid += middlePage
	}

	log.Printf("Sum of middle page numbers for valid updates: %d", sumValid)
	log.Printf("Sum of middle page numbers for invalid updates: %d", sumInvalid)
}
