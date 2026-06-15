package gotraverse

import "fmt"

// Algorithm is a graph search strategy. Implement it to add custom algorithms;
// the built-in strategies are [BFS], [DFS], [UCS], [Greedy] and [AStar].
type Algorithm interface {
	// Name returns a human-readable identifier used to populate Result.Algorithm.
	Name() string
	// Search finds a path from start to goal in g. Implementations must not
	// mutate g, so a single Graph can be searched repeatedly.
	Search(g *Graph, start, goal string) (Result, error)
}

// Search runs the given algorithm against this graph. It is a convenience
// wrapper equivalent to algo.Search(g, start, goal) that also validates that
// both endpoints exist.
func (g *Graph) Search(algo Algorithm, start, goal string) (Result, error) {
	if err := g.validate(start, goal); err != nil {
		return Result{Algorithm: algo.Name()}, err
	}
	return algo.Search(g, start, goal)
}

// validate checks that both endpoints exist in the graph.
func (g *Graph) validate(start, goal string) error {
	if !g.Has(start) {
		return fmt.Errorf("gotraverse: start node %q not in graph", start)
	}
	if !g.Has(goal) {
		return fmt.Errorf("gotraverse: goal node %q not in graph", goal)
	}
	return nil
}
