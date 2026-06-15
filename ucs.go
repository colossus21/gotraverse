package gotraverse

// UCS is uniform-cost search (Dijkstra from a single source). It expands nodes
// in order of accumulated path cost g(n) and returns a minimum-cost path,
// ignoring heuristics.
type UCS struct{}

func (UCS) Name() string { return "UCS" }

func (UCS) Search(g *Graph, start, goal string) (Result, error) {
	return bestFirst(g, start, goal, "UCS", true, func(_ string, gCost int) int {
		return gCost
	})
}
