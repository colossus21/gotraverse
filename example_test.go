package gotraverse_test

import (
	"fmt"

	"github.com/colossus21/gotraverse"
)

// ExampleParse builds a graph from the string format and runs A* on it.
func ExampleParse() {
	g, err := gotraverse.Parse(
		"S 8 A 8 B 4 C 3 D inf E inf G 0",
		"S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5",
	)
	if err != nil {
		panic(err)
	}

	res, _ := g.Search(gotraverse.AStar{}, "S", "G")
	fmt.Println(res.Path, res.Cost)
	// Output: [S C G] 13
}

// ExampleGraph_Search compares the path each algorithm returns on the same
// graph: the cost-aware searches find the cheaper S->C->G, while BFS returns
// the fewest-edge S->A->G.
func ExampleGraph_Search() {
	g, _ := gotraverse.Parse(
		"S 8 A 8 B 4 C 3 D inf E inf G 0",
		"S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5",
	)

	for _, algo := range []gotraverse.Algorithm{
		gotraverse.BFS{}, gotraverse.UCS{}, gotraverse.AStar{},
	} {
		res, _ := g.Search(algo, "S", "G")
		fmt.Printf("%-3s %v cost=%d\n", res.Algorithm, res.Path, res.Cost)
	}
	// Output:
	// BFS [S A G] cost=18
	// UCS [S C G] cost=13
	// A*  [S C G] cost=13
}

// ExampleNewGraph builds a graph programmatically instead of parsing strings.
func ExampleNewGraph() {
	g := gotraverse.NewGraph()
	g.AddNode("S", 1)
	g.AddNode("A", 1)
	g.AddNode("G", 0)
	_ = g.AddEdge("S", "A", 2)
	_ = g.AddEdge("A", "G", 5)

	res, _ := g.Search(gotraverse.UCS{}, "S", "G")
	fmt.Println(res.Path, res.Cost)
	// Output: [S A G] 7
}
