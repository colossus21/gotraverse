package gotraverse

// BFS is breadth-first search. It explores the graph level by level and ignores
// edge weights and heuristics, so the path it returns has the fewest edges but
// not necessarily the lowest cost.
func BFS[N comparable](p Problem[N]) (Result[N], error) {
	res := Result[N]{Algorithm: "BFS"}
	if err := p.validate(); err != nil {
		return res, err
	}

	dist := map[N]float64{p.Start: 0}
	visited := map[N]bool{p.Start: true}
	parent := map[N]N{}
	queue := []N{p.Start}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		res.Order = append(res.Order, cur)

		if p.Goal(cur) {
			res.Found = true
			res.Path = reconstructPath(parent, p.Start, cur)
			res.Cost = dist[cur]
			return res, nil
		}

		for _, e := range p.Neighbors(cur) {
			if !visited[e.To] {
				visited[e.To] = true
				parent[e.To] = cur
				dist[e.To] = dist[cur] + e.Weight
				queue = append(queue, e.To)
			}
		}
	}

	return res, nil
}
