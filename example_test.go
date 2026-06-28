package gotraverse_test

import (
	"context"
	"fmt"

	"github.com/colossus21/gotraverse"
)

// ExampleParse builds an explicit graph from the string format and runs A*.
func ExampleParse() {
	g, _ := gotraverse.Parse(
		"S 8 A 8 B 4 C 3 D inf E inf G 0",
		"S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5",
	)
	res, _ := gotraverse.AStar(g.Problem("S", "G"))
	fmt.Println(res.Path, res.Cost)
	// Output: [S C G] 13
}

// ExampleProblem searches an implicit graph that is never materialised: from n
// you may step to 2n or 2n+1. Neighbours are generated on demand, so the graph
// can be effectively unbounded.
func ExampleProblem() {
	p := gotraverse.Problem[int]{
		Start: 1,
		Goal:  gotraverse.GoalNode(5),
		Neighbors: func(n int) []gotraverse.Edge[int] {
			if n > 5 {
				return nil
			}
			return []gotraverse.Edge[int]{{To: 2 * n, Weight: 1}, {To: 2*n + 1, Weight: 1}}
		},
	}
	res, _ := gotraverse.AStar(p)
	fmt.Println(res.Path, res.Cost)
	// Output: [1 2 5] 2
}

// ExampleProblem_withContext shows cancelling a search. A cancelled context
// makes any algorithm abandon the search and return the context's error.
func ExampleProblem_withContext() {
	g, _ := gotraverse.Parse(
		"S 8 A 8 B 4 C 3 D inf E inf G 0",
		"S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5",
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // e.g. a timeout fired or the caller gave up

	_, err := gotraverse.AStar(g.Problem("S", "G").WithContext(ctx))
	fmt.Println(err)
	// Output: context canceled
}

// ExampleDepthLimited shows the configurable depth cutoff: G sits two edges from
// S, so a limit of 1 cannot reach it but a limit of 2 can.
func ExampleDepthLimited() {
	g, _ := gotraverse.Parse(
		"S 8 A 8 B 4 C 3 D inf E inf G 0",
		"S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5",
	)
	shallow, _ := gotraverse.DepthLimited[string](1)(g.Problem("S", "G"))
	deep, _ := gotraverse.DepthLimited[string](2)(g.Problem("S", "G"))
	fmt.Println(shallow.Found)
	fmt.Println(deep.Found, deep.Path)
	// Output:
	// false
	// true [S A G]
}

// Example_strategies runs several algorithms over the same problem by treating
// them as interchangeable SearchFunc values.
func Example_strategies() {
	g, _ := gotraverse.Parse(
		"S 8 A 8 B 4 C 3 D inf E inf G 0",
		"S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5",
	)
	p := g.Problem("S", "G")
	for _, search := range []gotraverse.SearchFunc[string]{
		gotraverse.BFS[string], gotraverse.UCS[string], gotraverse.AStar[string],
	} {
		res, _ := search(p)
		fmt.Printf("%-3s %v cost=%v\n", res.Algorithm, res.Path, res.Cost)
	}
	// Output:
	// BFS [S A G] cost=18
	// UCS [S C G] cost=13
	// A*  [S C G] cost=13
}
