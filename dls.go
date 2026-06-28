package gotraverse

// DepthLimited returns depth-first search bounded to a maximum depth (number of
// edges from the start). It is the building block of [IDDFS] and is useful on
// its own to cap how deep a search may go. Like DFS it ignores edge weights and
// heuristics.
//
// Depth is measured in edges, so limit 0 explores only the start node, limit 1
// reaches its direct successors, and so on. Neighbours are visited in the order
// Neighbors returns them, and cycles are avoided along the current path so the
// search always terminates.
func DepthLimited[N comparable](limit int) SearchFunc[N] {
	return func(p Problem[N]) (Result[N], error) {
		res := Result[N]{Algorithm: "Depth-Limited"}
		if err := p.validate(); err != nil {
			return res, err
		}
		r, _ := depthLimited(p, limit, "Depth-Limited")
		if err := p.cancelled(); err != nil {
			return r, err
		}
		return r, nil
	}
}

// depthLimited runs a single depth-limited DFS and reports whether a cutoff
// occurred — i.e. some node was left unexpanded only because the depth limit was
// reached. A pass with no cutoff means the reachable space was fully explored,
// which lets [IDDFS] stop without knowing the graph size.
func depthLimited[N comparable](p Problem[N], limit int, name string) (Result[N], bool) {
	res := Result[N]{Algorithm: name}
	onPath := map[N]bool{}
	path := []N{}
	cutoff := false
	aborted := false

	var dfs func(node N, depth int, cost float64) bool
	dfs = func(node N, depth int, cost float64) bool {
		if aborted {
			return false
		}
		if p.cancelled() != nil {
			aborted = true
			return false // unwind; the caller reports ctx.Err()
		}
		res.Order = append(res.Order, node)
		path = append(path, node)
		onPath[node] = true
		defer func() {
			path = path[:len(path)-1]
			onPath[node] = false
		}()

		if p.Goal(node) {
			res.Found = true
			res.Path = append([]N(nil), path...)
			res.Cost = cost
			return true
		}
		if depth >= limit {
			if len(p.Neighbors(node)) > 0 {
				cutoff = true // there is more graph below this node
			}
			return false
		}
		for _, e := range p.Neighbors(node) {
			if onPath[e.To] {
				continue // skip nodes already on the current path (cycle)
			}
			if dfs(e.To, depth+1, cost+e.Weight) {
				return true
			}
		}
		return false
	}

	dfs(p.Start, 0, 0)
	return res, cutoff
}
