package gotraverse

// BFS is breadth-first search. It explores the graph level by level and ignores
// edge weights and heuristics, so the path it returns has the fewest edges but
// not necessarily the lowest cost.
type BFS struct{}

func (BFS) Name() string { return "BFS" }

func (BFS) Search(g *Graph, start, goal string) (Result, error) {
	res := Result{Algorithm: "BFS"}
	if err := g.validate(start, goal); err != nil {
		return res, err
	}

	visited := map[string]bool{start: true}
	parent := map[string]string{}
	queue := []string{start}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		res.Order = append(res.Order, cur)

		if cur == goal {
			res.Found = true
			res.Path = reconstructPath(parent, start, goal)
			res.Cost = pathCost(g, res.Path)
			return res, nil
		}

		for _, e := range g.adj[cur] {
			if !visited[e.to] {
				visited[e.to] = true
				parent[e.to] = cur
				queue = append(queue, e.to)
			}
		}
	}

	return res, nil
}
