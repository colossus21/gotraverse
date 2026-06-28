package gotraverse

// Result is the outcome of a search.
type Result[N comparable] struct {
	// Algorithm is the human-readable name of the algorithm that produced this
	// result (e.g. "A*").
	Algorithm string
	// Found reports whether a goal was reached from the start node.
	Found bool
	// Path is the sequence of nodes from start to goal, inclusive. It is nil
	// when Found is false.
	Path []N
	// Cost is the total weight of the edges along Path. For unweighted searches
	// (BFS, DFS and friends) it is still the summed edge weight of the path
	// actually returned, which need not be minimal.
	Cost float64
	// Order is the sequence in which nodes were expanded, useful for tracing.
	// Iterative algorithms (IDDFS, IDA*) repeat nodes across their passes.
	Order []N
}

// reconstructPath walks the parent map back from goal to start and returns the
// forward path. start must be reachable from goal via parent links.
func reconstructPath[N comparable](parent map[N]N, start, goal N) []N {
	path := []N{goal}
	for cur := goal; cur != start; {
		p := parent[cur]
		path = append(path, p)
		cur = p
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}
