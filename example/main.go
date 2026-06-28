// Command example demonstrates gotraverse on both an explicit string graph and
// an implicit 2D grid solved with A* — the implicit case never materialises the
// whole graph, generating neighbours on demand.
package main

import (
	"fmt"
	"log"
	"math"

	"github.com/colossus21/gotraverse"
)

func main() {
	explicitGraph()
	fmt.Println()
	gridPathfinding()
}

func explicitGraph() {
	g, err := gotraverse.Parse(
		"S 8 A 8 B 4 C 3 D inf E inf G 0",
		"S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5",
	)
	if err != nil {
		log.Fatal(err)
	}
	p := g.Problem("S", "G")

	searches := []gotraverse.SearchFunc[string]{
		gotraverse.BFS[string],
		gotraverse.DFS[string],
		gotraverse.DepthLimited[string](5),
		gotraverse.IDDFS[string],
		gotraverse.Bidirectional[string],
		gotraverse.UCS[string],
		gotraverse.Greedy[string],
		gotraverse.AStar[string],
		gotraverse.IDAStar[string],
	}

	fmt.Println("# Explicit graph: S -> G")
	for _, search := range searches {
		res, err := search(p)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-14s %v  cost=%v\n", res.Algorithm, res.Path, res.Cost)
	}
}

// Point is a grid cell. Used as the node type of an implicit graph.
type Point struct{ X, Y int }

func gridPathfinding() {
	grid := []string{
		"S....",
		".###.",
		"...#.",
		".#.#.",
		".#..G",
	}

	var start, goal Point
	for y, row := range grid {
		for x := range row {
			switch row[x] {
			case 'S':
				start = Point{x, y}
			case 'G':
				goal = Point{x, y}
			}
		}
	}

	walkable := func(p Point) bool {
		return p.Y >= 0 && p.Y < len(grid) &&
			p.X >= 0 && p.X < len(grid[p.Y]) &&
			grid[p.Y][p.X] != '#'
	}

	// The graph is never built: neighbours are generated on the fly.
	p := gotraverse.Problem[Point]{
		Start: start,
		Goal:  gotraverse.GoalNode(goal),
		Neighbors: func(p Point) []gotraverse.Edge[Point] {
			var es []gotraverse.Edge[Point]
			for _, d := range []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
				n := Point{p.X + d.X, p.Y + d.Y}
				if walkable(n) {
					es = append(es, gotraverse.Edge[Point]{To: n, Weight: 1})
				}
			}
			return es
		},
		Heuristic: func(p Point) float64 {
			return math.Abs(float64(p.X-goal.X)) + math.Abs(float64(p.Y-goal.Y))
		},
	}

	res, err := gotraverse.AStar(p)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("# Implicit grid: A* from S to G")
	if !res.Found {
		fmt.Println("no path")
		return
	}

	overlay := make([][]byte, len(grid))
	for y := range grid {
		overlay[y] = []byte(grid[y])
	}
	for _, pt := range res.Path {
		if c := overlay[pt.Y][pt.X]; c != 'S' && c != 'G' {
			overlay[pt.Y][pt.X] = '*'
		}
	}
	for _, row := range overlay {
		fmt.Println(string(row))
	}
	fmt.Printf("path length: %v steps\n", res.Cost)
}
