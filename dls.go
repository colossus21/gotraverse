package gotraverse

// DepthLimited is depth-first search bounded to a maximum depth (number of
// edges from the start). It is the building block of [IterativeDeepening] and
// is useful on its own to cap how deep a search may go. Like DFS it ignores
// edge weights and heuristics.
//
// Depth is measured in edges, so Limit 0 explores only the start node, Limit 1
// reaches its direct successors, and so on. Children are visited in edge
// declaration order, and cycles are avoided along the current path so the
// search always terminates.
type DepthLimited struct {
	// Limit is the maximum depth (in edges) to explore.
	Limit int
}

func (DepthLimited) Name() string { return "Depth-Limited" }

func (d DepthLimited) Search(g *Graph, start, goal string) (Result, error) {
	res := Result{Algorithm: d.Name()}
	if err := g.validate(start, goal); err != nil {
		return res, err
	}

	onPath := map[string]bool{}
	path := []string{}

	var dfs func(node string, depth, cost int) bool
	dfs = func(node string, depth, cost int) bool {
		res.Order = append(res.Order, node)
		path = append(path, node)
		onPath[node] = true
		defer func() {
			path = path[:len(path)-1]
			onPath[node] = false
		}()

		if node == goal {
			res.Found = true
			res.Path = append([]string(nil), path...)
			res.Cost = cost
			return true
		}
		if depth >= d.Limit {
			return false // depth limit reached: cutoff
		}
		for _, e := range g.adj[node] {
			if onPath[e.to] {
				continue // skip nodes already on the current path (cycle)
			}
			if dfs(e.to, depth+1, cost+e.weight) {
				return true
			}
		}
		return false
	}

	dfs(start, 0, 0)
	return res, nil
}
