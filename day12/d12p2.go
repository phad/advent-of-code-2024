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
				if cellAdjoins(n1.cell, n2.cell) {
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

func cellAdjoins(r1, r2 aoc.Point) bool {
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

func (e edge) String() string {
	return map[edge]string{
		top:    "top",
		right:  "right",
		bottom: "bottom",
		left:   "left",
	}[e]
}

type winding int

const (
	cw winding = iota
	ccw
)

func (w winding) String() string {
	return map[winding]string{
		cw:  "cw",
		ccw: "ccw",
	}[w]
}

type fence struct {
	cell aoc.Point
	edge edge
}

func (f fence) String() string {
	return fmt.Sprintf("[%v %v]", f.cell, f.edge)
}

func (r *region) findPanels() map[fence]map[winding]int {
	panels := map[fence]map[winding]int{}
	inc := func(c aoc.Point, e edge, w winding) {
		f := fence{cell: c, edge: e}
		wc, ok := panels[f]
		if !ok {
			wc = map[winding]int{}
			panels[f] = wc
		}
		panels[f][w]++
	}

	for _, c := range r.cells {
		inc(c, top, cw)
		inc(c, right, cw)
		inc(aoc.Point{c.X, c.Y + 1}, top, ccw /*c bottom cw*/)
		inc(aoc.Point{c.X - 1, c.Y}, right, ccw /*c left cw*/)
	}
	return panels
}

func (r *region) perimeter() int {
	ret := []fence{}
	for f, windingCounts := range r.findPanels() {
		if len(windingCounts) == 1 {
			ret = append(ret, f)
		}
	}
	return len(ret)
}

type wFence struct {
	f fence
	w winding
}

func (r *region) woundFences() []wFence {
	var ret []wFence
	for f, windingCounts := range r.findPanels() {
		if len(windingCounts) == 1 {
			if _, ok := windingCounts[cw]; ok {
				ret = append(ret, wFence{f: f, w: cw})
			} else {
				ret = append(ret, wFence{f: f, w: ccw})
			}
		}
	}
	return ret
}

func (r *region) sides() int {
	parents := map[wFence]wFence{}
	fences := r.woundFences()
	for _, f := range fences {
		parents[f] = f
	}
	//log.Printf("fences: %v\nparents: %v\n", fences, parents)

	rootFn := func(f wFence) wFence {
		var p, root wFence
		p = parents[f]
		for {
			if parents[p] == p {
				root = p
				break
			}
			p = parents[p]
		}
		return root
	}
	for i := 0; i < len(fences)-1; i++ {
		for j := 0; j < len(fences); j++ {
			f1, f2 := fences[i], fences[j]
			if fenceAdjoins(f1, f2) {
				parents[rootFn(f1)] = rootFn(f2)
			}
		}
	}
	//log.Printf("Parents: %v", parents)
	sides := map[wFence]int{}
	for _, f := range fences {
		r := rootFn(f)
		sides[r]++
	}
	//log.Printf("Sides: %v", sides)
	return len(sides)
}

func fenceAdjoins(a, b wFence) bool {
	if a.f.edge != b.f.edge {
		return false
	}
	if a.w != b.w {
		return false
	}
	switch a.f.edge {
	case top:
		if a.f.cell.Y != b.f.cell.Y {
			return false
		}
		if abs(b.f.cell.X-a.f.cell.X) != 1 {
			return false
		}
	case right:
		if a.f.cell.X != b.f.cell.X {
			return false
		}
		if abs(b.f.cell.Y-a.f.cell.Y) != 1 {
			return false
		}
	}
	return true
}

func main() {
	log.Println("AoC-2024-day12-part2")
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
		sides := reg.sides()
		cost := area * sides
		totalCost += cost
		log.Printf("Plant %s:\n%vArea: %d\nPerimeter: %d\nSides: %d\nCost: %d\n\n", string(reg.plant), g.Highlight(reg.plant), area, perim, sides, cost)
	}
	log.Printf("Total cost: %d", totalCost)
}
