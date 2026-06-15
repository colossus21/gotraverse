package gotraverse

// DFS is depth-first search. It dives along each branch before backtracking and
// ignores edge weights and heuristics. Neighbours are pushed in declaration
// order, so the last-declared edge of a node is explored first.
type DFS struct{}

func (DFS) Name() string { return "DFS" }

// frame is a stack entry pairing a node with the node that pushed it, so the
// DFS-tree parent is exactly the node from which it was actually reached.
type frame struct{ node, parent string }

func (DFS) Search(g *Graph, start, goal string) (Result, error) {
	res := Result{Algorithm: "DFS"}
	if err := g.validate(start, goal); err != nil {
		return res, err
	}

	visited := map[string]bool{}
	parent := map[string]string{}
	stack := []frame{{node: start, parent: ""}}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if visited[top.node] {
			continue
		}
		visited[top.node] = true
		if top.node != start {
			parent[top.node] = top.parent
		}
		res.Order = append(res.Order, top.node)

		if top.node == goal {
			res.Found = true
			res.Path = reconstructPath(parent, start, goal)
			res.Cost = pathCost(g, res.Path)
			return res, nil
		}

		for _, e := range g.adj[top.node] {
			if !visited[e.to] {
				stack = append(stack, frame{node: e.to, parent: top.node})
			}
		}
	}

	return res, nil
}
