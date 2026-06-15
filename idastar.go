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
type IDAStar struct{}

func (IDAStar) Name() string { return "IDA*" }

func (IDAStar) Search(g *Graph, start, goal string) (Result, error) {
	res := Result{Algorithm: "IDA*"}
	if err := g.validate(start, goal); err != nil {
		return res, err
	}

	h := func(n string) int {
		v, _ := g.Heuristic(n)
		return v
	}

	onPath := map[string]bool{}
	path := []string{}

	// dfs returns whether the goal was found, and the smallest f value that
	// exceeded the current threshold (math.MaxInt if none did), which becomes
	// the next pass's threshold.
	var dfs func(node string, gCost, threshold int) (bool, int)
	dfs = func(node string, gCost, threshold int) (bool, int) {
		f := gCost + h(node)
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

		if node == goal {
			res.Found = true
			res.Path = append([]string(nil), path...)
			res.Cost = gCost
			return true, threshold
		}

		next := math.MaxInt
		for _, e := range g.adj[node] {
			if onPath[e.to] {
				continue
			}
			found, t := dfs(e.to, gCost+e.weight, threshold)
			if found {
				return true, t
			}
			if t < next {
				next = t
			}
		}
		return false, next
	}

	threshold := h(start)
	for {
		found, t := dfs(start, 0, threshold)
		if found {
			return res, nil
		}
		if t == math.MaxInt {
			return res, nil // exhausted the reachable space: goal unreachable
		}
		threshold = t
	}
}
