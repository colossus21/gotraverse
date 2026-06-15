package gotraverse

// AStar is A* search. It expands nodes in order of f(n) = g(n) + h(n), the
// accumulated path cost plus the heuristic estimate to the goal. With an
// admissible (non-overestimating) heuristic it returns a minimum-cost path.
type AStar struct{}

func (AStar) Name() string { return "A*" }

func (AStar) Search(g *Graph, start, goal string) (Result, error) {
	return bestFirst(g, start, goal, "A*", true, func(node string, gCost int) int {
		h, _ := g.Heuristic(node)
		// Guard against overflow when h is Inf (unreachable-by-heuristic nodes).
		if h >= Inf-gCost {
			return Inf
		}
		return gCost + h
	})
}
