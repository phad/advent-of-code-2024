package main

import (
	"fmt"
	"log"
	"os"
)

/* Example input
AAAA
BBCD
BBCC
EEEC
*/

type grid struct {
	w, h  int
	cells [][]rune
}

func newGrid(in []string) (*grid, error) {
	g := &grid{h: len(in)}
	for i, r := range in {
		if i == 0 {
			g.w = len(r)
			if g.w != g.h {
				return nil, fmt.Errorf("Grid isn't square: width %d != height %d", g.w, g.h)
			}
		}
		if i > 0 && len(r) != g.w {
			return nil, fmt.Errorf("Row %d wrong size %d want %d", 1, len(r), g.w)
		}
		row := []rune(r)
		g.cells = append(g.cells, row)
	}
	return g, nil
}

func (g *grid) String() string {
	s := fmt.Sprintf("width:%d height:%d\n", g.w, g.h)
	for _, r := range g.cells {
		s += string(r)
		s += "\n"
	}
	return s
}

func (g *grid) highlight(show rune) string {
	s := fmt.Sprintf("width:%d height:%d\n", g.w, g.h)
	for _, row := range g.cells {
		r := make([]rune, g.w)
		copy(r, row)
		for i, c := range row {
			if show != c {
				r[i] = '.'
			}
		}
		s += string(r)
		s += "\n"
	}
	return s
}

type point struct{ x, y int }

type region struct {
	plant rune
	cells []point
}

func (r *region) String() string {
	return fmt.Sprintf("<%s: %v>", string(r.plant), r.cells)
}

type node struct {
	cell point
}

func (n *node) String() string {
	return fmt.Sprintf("<%v>", n.cell)
}

func (g *grid) findRegions() []*region {
	// Start by creating a lot of 1-cell nodes for union-find.
	cellsByPlant := map[rune][]*node{}
	for y, row := range g.cells {
		for x, plant := range row {
			cs, ok := cellsByPlant[plant]
			if !ok {
				cs = []*node{}
				cellsByPlant[plant] = cs
			}
			n := &node{cell: point{x, y}}
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

func adjoins(r1, r2 point) bool {
	if r1.x == r2.x {
		return abs(r2.y-r1.y) == 1
	}
	if r1.y == r2.y {
		return abs(r2.x-r1.x) == 1
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
	cell point
	edge edge
}

func (r *region) findPanels() map[fence]int {
	panels := map[fence]int{}
	for _, c := range r.cells {
		panels[fence{cell: c, edge: top}]++
		panels[fence{cell: c, edge: right}]++
		panels[fence{cell: point{c.x, c.y + 1}, edge: top /*c bottom*/}]++
		panels[fence{cell: point{c.x - 1, c.y}, edge: right /*c left*/}]++
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
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	g, err := newGrid(lines)
	if err != nil {
		log.Fatalf("Failed to parse input: %v", err)
	}
	log.Printf("AllPlants:\n%v\n", g)

	totalCost := 0
	for _, reg := range g.findRegions() {
		area := reg.area()
		perim := reg.perimeter()
		cost := area * perim
		totalCost += cost
		log.Printf("Plant %s:\n%vArea: %d\nPerimeter: %d\nCost: %d\n\n", string(reg.plant), g.highlight(reg.plant), area, perim, cost)
	}
	log.Printf("Total cost: %d", totalCost)
}
