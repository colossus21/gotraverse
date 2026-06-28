package gotraverse

// IDDFS is iterative-deepening depth-first search. It runs depth-limited DFS
// repeatedly with an increasing limit (0, 1, 2, …) until the goal is found,
// combining DFS's low memory use with BFS's guarantee of finding the shallowest
// goal (fewest edges). Edge weights and heuristics are ignored.
//
// It stops when a pass finds the goal or when a pass completes without any
// depth cutoff — meaning the reachable space is finite and fully explored — so
// it terminates on finite graphs without needing to know their size. On an
// infinite graph with no reachable goal it does not terminate, as is inherent
// to the algorithm. Order accumulates the expansions of every pass, so earlier
// nodes reappear on each iteration.
func IDDFS[N comparable](p Problem[N]) (Result[N], error) {
	res := Result[N]{Algorithm: "IDDFS"}
	if err := p.validate(); err != nil {
		return res, err
	}

	for limit := 0; ; limit++ {
		sub, cutoff := depthLimited(p, limit, "IDDFS")
		res.Order = append(res.Order, sub.Order...)
		if err := p.cancelled(); err != nil {
			return res, err
		}
		if sub.Found {
			res.Found = true
			res.Path = sub.Path
			res.Cost = sub.Cost
			return res, nil
		}
		if !cutoff {
			return res, nil // reachable space exhausted: goal unreachable
		}
	}
}
