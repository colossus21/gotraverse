package gotraverse

// IterativeDeepening is iterative-deepening depth-first search (IDDFS). It runs
// [DepthLimited] repeatedly with an increasing limit (0, 1, 2, …) until the
// goal is found, combining DFS's low memory use with BFS's guarantee of finding
// the shallowest goal (fewest edges). Edge weights and heuristics are ignored.
//
// Order accumulates the expansions of every deepening pass, so earlier nodes
// reappear on each iteration — that repetition is inherent to IDDFS.
type IterativeDeepening struct{}

func (IterativeDeepening) Name() string { return "IDDFS" }

func (IterativeDeepening) Search(g *Graph, start, goal string) (Result, error) {
	res := Result{Algorithm: "IDDFS"}
	if err := g.validate(start, goal); err != nil {
		return res, err
	}

	// The deepest a simple (cycle-free) path can be is len(nodes)-1 edges;
	// beyond that the goal is unreachable.
	maxDepth := len(g.heuristic) - 1
	for limit := 0; limit <= maxDepth; limit++ {
		sub, _ := DepthLimited{Limit: limit}.Search(g, start, goal)
		res.Order = append(res.Order, sub.Order...)
		if sub.Found {
			res.Found = true
			res.Path = sub.Path
			res.Cost = sub.Cost
			return res, nil
		}
	}
	return res, nil
}
