package gotraverse

// DFS is depth-first search. It dives along each branch before backtracking and
// ignores edge weights and heuristics. Neighbours are pushed in the order
// Neighbors returns them, so the last is explored first.
func DFS[N comparable](p Problem[N]) (Result[N], error) {
	res := Result[N]{Algorithm: "DFS"}
	if err := p.validate(); err != nil {
		return res, err
	}

	type frame[T comparable] struct {
		node    T
		parent  T
		hasPrev bool
		cost    float64
	}

	visited := map[N]bool{}
	parent := map[N]N{}
	dist := map[N]float64{}
	stack := []frame[N]{{node: p.Start}}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if visited[top.node] {
			continue
		}
		visited[top.node] = true
		if top.hasPrev {
			parent[top.node] = top.parent
		}
		dist[top.node] = top.cost
		res.Order = append(res.Order, top.node)

		if p.Goal(top.node) {
			res.Found = true
			res.Path = reconstructPath(parent, p.Start, top.node)
			res.Cost = top.cost
			return res, nil
		}

		for _, e := range p.Neighbors(top.node) {
			if !visited[e.To] {
				stack = append(stack, frame[N]{
					node: e.To, parent: top.node, hasPrev: true, cost: top.cost + e.Weight,
				})
			}
		}
	}

	return res, nil
}
