package gotraverse

import "sort"

// Bidirectional is bidirectional breadth-first search. It grows two frontiers
// at once — forward from the start along out-edges and backward from the goal
// along in-edges — always expanding the smaller frontier, until they meet. The
// path it stitches together has the fewest edges, like [BFS], but it typically
// touches far fewer nodes. Edge weights and heuristics are ignored.
type Bidirectional struct{}

func (Bidirectional) Name() string { return "Bidirectional" }

func (Bidirectional) Search(g *Graph, start, goal string) (Result, error) {
	res := Result{Algorithm: "Bidirectional"}
	if err := g.validate(start, goal); err != nil {
		return res, err
	}

	if start == goal {
		res.Found = true
		res.Path = []string{start}
		res.Order = []string{start}
		return res, nil
	}

	// Reverse adjacency for the backward search. Built by iterating source
	// nodes in sorted order so predecessor lists — and therefore the meeting
	// node and resulting path — are deterministic.
	reverse := make(map[string][]string)
	sources := make([]string, 0, len(g.adj))
	for u := range g.adj {
		sources = append(sources, u)
	}
	sort.Strings(sources)
	for _, u := range sources {
		for _, e := range g.adj[u] {
			reverse[e.to] = append(reverse[e.to], u)
		}
	}

	parentF := map[string]string{start: ""} // node -> predecessor toward start
	parentB := map[string]string{goal: ""}  // node -> successor toward goal
	frontierF := []string{start}
	frontierB := []string{goal}

	fwd := func(n string) []string {
		outs := make([]string, 0, len(g.adj[n]))
		for _, e := range g.adj[n] {
			outs = append(outs, e.to)
		}
		return outs
	}
	bwd := func(n string) []string { return reverse[n] }

	// expand advances one frontier by a single BFS layer. It records each new
	// node's parent and returns the next layer plus a meeting node (if any
	// neighbour was already reached by the opposite search).
	expand := func(frontier []string, parent, other map[string]string, adj func(string) []string) ([]string, string) {
		var next []string
		for _, node := range frontier {
			res.Order = append(res.Order, node)
			for _, nb := range adj(node) {
				if _, seen := parent[nb]; seen {
					continue
				}
				parent[nb] = node
				if _, met := other[nb]; met {
					return next, nb
				}
				next = append(next, nb)
			}
		}
		return next, ""
	}

	meet := ""
	for len(frontierF) > 0 && len(frontierB) > 0 {
		var m string
		if len(frontierF) <= len(frontierB) {
			frontierF, m = expand(frontierF, parentF, parentB, fwd)
		} else {
			frontierB, m = expand(frontierB, parentB, parentF, bwd)
		}
		if m != "" {
			meet = m
			break
		}
	}

	if meet == "" {
		return res, nil // frontiers never met: goal unreachable
	}

	// Build start..meet from the forward tree, then meet..goal from the
	// backward tree.
	var left []string
	for n := meet; n != ""; n = parentF[n] {
		left = append(left, n)
	}
	for i, j := 0, len(left)-1; i < j; i, j = i+1, j-1 {
		left[i], left[j] = left[j], left[i]
	}
	path := left
	for n := parentB[meet]; n != ""; n = parentB[n] {
		path = append(path, n)
	}

	res.Found = true
	res.Path = path
	res.Cost = pathCost(g, path)
	return res, nil
}
