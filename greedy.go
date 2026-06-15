package gotraverse

// Greedy is greedy best-first search. It always expands the open node with the
// smallest heuristic h(n), ignoring the cost already spent. It is fast but not
// optimal: the returned path is the first one found to the goal.
type Greedy struct{}

func (Greedy) Name() string { return "Greedy" }

func (Greedy) Search(g *Graph, start, goal string) (Result, error) {
	return bestFirst(g, start, goal, "Greedy", false, func(node string, _ int) int {
		h, _ := g.Heuristic(node)
		return h
	})
}
