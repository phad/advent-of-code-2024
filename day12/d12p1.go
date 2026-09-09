package main

import (
	"fmt"
	"log"
	"os"

	"github.com/phad/advent-of-code-2024/aoc"
)

/* Example input
AAAA
BBCD
BBCC
EEEC
*/

type region struct {
	plant rune
	cells []aoc.Point
}

func (r *region) String() string {
	return fmt.Sprintf("<%s: %v>", string(r.plant), r.cells)
}

type node struct {
	cell aoc.Point
}

func (n *node) String() string {
	return fmt.Sprintf("<%v>", n.cell)
}

func findRegions(g *aoc.Grid[rune]) []*region {
	// Start by creating a lot of 1-cell nodes for union-find.
	cellsByPlant := map[rune][]*node{}
	for y, row := range g.Cells {
		for x, plant := range row {
			cs, ok := cellsByPlant[plant]
			if !ok {
				cs = []*node{}
				cellsByPlant[plant] = cs
			}
			n := &node{cell: aoc.Point{x, y}}
			cellsByPlant[plant] = append(cellsByPlant[plant], n)
		}
	}

	var ret []*region

	for plant, nodes := range cellsByPlant {
		// Now use union-find to merge cells into regions, where
		// all cells adjoin on 1 or more sides.  Do this per plant
		// so that we end up with disjoint regions for a particular
		// plant.  Initially every node is it's own parent.
		parents := map[*node]*node{}
		for _, n := range nodes {
			parents[n] = n
		}
		//log.Printf("Plant %s: initial parents:\n%v", string(plant), parents)

		rootFn := func(n *node) *node {
			var p, root *node
			p = parents[n]
			for {
				if parents[p] == p {
					root = p
					break
				}
				p = parents[p]
			}
			return root
		}

		// Union
		for i := 0; i < len(nodes)-1; i++ {
			for j := i + 1; j < len(nodes); j++ {
				n1, n2 := nodes[i], nodes[j]
				if adjoins(n1.cell, n2.cell) {
					parents[rootFn(n2)] = rootFn(n1)
				}
			}
		}
		// Find
		//log.Printf("Plant %s: union-find parents state:\n%v\n", string(plant), parents)

		// Create output regions - need to map each cluster's root to the new region.
		// 1x1 islands don't have a root in the parents list.
		regions := map[*node]*region{}
		for _, n := range nodes {
			root := rootFn(n)
			//log.Printf("For node %v found root %v", n, root)
			reg, ok := regions[root]
			if !ok {
				reg = &region{plant: plant}
				regions[root] = reg
			}
			//log.Printf("For root %v found region %v", root, reg)
			reg.cells = append(reg.cells, n.cell)
		}
		//log.Printf("Made regions:\n%v", regions)
		for _, r := range regions {
			ret = append(ret, r)
		}
	}
	return ret
}

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func adjoins(r1, r2 aoc.Point) bool {
	if r1.X == r2.X {
		return abs(r2.Y-r1.Y) == 1
	}
	if r1.Y == r2.Y {
		return abs(r2.X-r1.X) == 1
	}
	return false
}

func (r *region) area() int {
	return len(r.cells)
}

type edge int

const (
	top edge = iota
	right
	bottom
	left
)

type fence struct {
	cell aoc.Point
	edge edge
}

func (r *region) findPanels() map[fence]int {
	panels := map[fence]int{}
	for _, c := range r.cells {
		panels[fence{cell: c, edge: top}]++
		panels[fence{cell: c, edge: right}]++
		panels[fence{cell: aoc.Point{c.X, c.Y + 1}, edge: top /*c bottom*/}]++
		panels[fence{cell: aoc.Point{c.X - 1, c.Y}, edge: right /*c left*/}]++
	}
	return panels
}

func (r *region) perimeter() int {
	singlePanels := 0
	for _, count := range r.findPanels() {
		if count == 1 {
			singlePanels++
		}
	}
	return singlePanels
}

func main() {
	log.Println("AoC-2024-day12-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := aoc.ReadLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	g, err := aoc.NewRuneGrid(lines, true /*=wantSquare*/)
	if err != nil {
		log.Fatalf("Failed to parse input: %v", err)
	}
	log.Printf("AllPlants:\n%v\n", g)

	totalCost := 0
	for _, reg := range findRegions(g) {
		area := reg.area()
		perim := reg.perimeter()
		cost := area * perim
		totalCost += cost
		log.Printf("Plant %s:\n%vArea: %d\nPerimeter: %d\nCost: %d\n\n", string(reg.plant), g.Highlight(reg.plant), area, perim, cost)
	}
	log.Printf("Total cost: %d", totalCost)
}
