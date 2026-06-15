package gotraverse

// Result is the outcome of a search.
type Result struct {
	// Algorithm is the human-readable name of the algorithm that produced this
	// result (e.g. "A*").
	Algorithm string
	// Found reports whether the goal was reached from the start node.
	Found bool
	// Path is the sequence of node names from start to goal, inclusive. It is
	// nil when Found is false.
	Path []string
	// Cost is the total weight of the edges along Path. For unweighted searches
	// (BFS, DFS) this is still the summed edge weight of the path actually
	// returned, which need not be minimal.
	Cost int
	// Order is the sequence in which nodes were expanded (dequeued/popped),
	// useful for tracing or visualising how a search proceeded.
	Order []string
}

// reconstructPath walks the parent map back from goal to start and returns the
// forward path. It assumes goal is reachable (present as a key or equal to
// start).
func reconstructPath(parent map[string]string, start, goal string) []string {
	path := []string{goal}
	for cur := goal; cur != start; {
		p := parent[cur]
		path = append(path, p)
		cur = p
	}
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// pathCost sums the edge weights along path. If two consecutive nodes are not
// directly connected (should not happen for a valid path) the missing hop
// contributes zero.
func pathCost(g *Graph, path []string) int {
	cost := 0
	for i := 0; i+1 < len(path); i++ {
		for _, e := range g.adj[path[i]] {
			if e.to == path[i+1] {
				cost += e.weight
				break
			}
		}
	}
	return cost
}
