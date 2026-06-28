package gotraverse

import "math"

// IDAStar is iterative-deepening A*. It performs a series of depth-first
// searches bounded by an f = g(n) + h(n) threshold rather than by depth,
// raising the threshold to the smallest f that exceeded it on each pass. With
// an admissible heuristic it returns a minimum-cost path using far less memory
// than [AStar] (no open/closed sets), at the cost of re-expanding nodes.
//
// Like A* and UCS it assumes non-negative edge weights. Order accumulates the
// expansions of every threshold pass.
func IDAStar[N comparable](p Problem[N]) (Result[N], error) {
	res := Result[N]{Algorithm: "IDA*"}
	if err := p.validate(); err != nil {
		return res, err
	}

	onPath := map[N]bool{}
	path := []N{}

	// dfs returns whether the goal was found, and the smallest f value that
	// exceeded the current threshold (+Inf if none did), which becomes the next
	// pass's threshold.
	var dfs func(node N, gCost, threshold float64) (bool, float64)
	dfs = func(node N, gCost, threshold float64) (bool, float64) {
		f := gCost + p.h(node)
		if f > threshold {
			return false, f
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
			res.Cost = gCost
			return true, threshold
		}

		next := math.Inf(1)
		for _, e := range p.Neighbors(node) {
			if onPath[e.To] {
				continue
			}
			found, t := dfs(e.To, gCost+e.Weight, threshold)
			if found {
				return true, t
			}
			if t < next {
				next = t
			}
		}
		return false, next
	}

	threshold := p.h(p.Start)
	for {
		found, t := dfs(p.Start, 0, threshold)
		if found {
			return res, nil
		}
		if math.IsInf(t, 1) {
			return res, nil // reachable space exhausted: goal unreachable
		}
		threshold = t
	}
}
