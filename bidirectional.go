package gotraverse

import "errors"

// Bidirectional is bidirectional breadth-first search. It grows two frontiers at
// once — forward from the start along out-edges and backward from the goal
// node(s) along in-edges — always expanding the smaller frontier, until they
// meet. The stitched path has the fewest edges, like [BFS], but it typically
// touches far fewer nodes. Edge weights and heuristics are ignored.
//
// Because it searches backward, it requires Problem.Predecessors and at least
// one Problem.GoalNodes entry; both are populated automatically by
// [Graph.Problem]. It returns an error if either is missing.
func Bidirectional[N comparable](p Problem[N]) (Result[N], error) {
	res := Result[N]{Algorithm: "Bidirectional"}
	if err := p.validate(); err != nil {
		return res, err
	}
	if p.Predecessors == nil {
		return res, errors.New("gotraverse: Bidirectional requires Problem.Predecessors")
	}
	if len(p.GoalNodes) == 0 {
		return res, errors.New("gotraverse: Bidirectional requires Problem.GoalNodes")
	}

	distF := map[N]float64{p.Start: 0}
	distB := map[N]float64{}
	parentF := map[N]N{}
	parentB := map[N]N{}
	rootF := map[N]bool{p.Start: true}
	rootB := map[N]bool{}

	frontierF := []N{p.Start}
	var frontierB []N
	for _, gn := range p.GoalNodes {
		if _, ok := distB[gn]; ok {
			continue
		}
		distB[gn] = 0
		rootB[gn] = true
		frontierB = append(frontierB, gn)
	}

	// Start is itself a goal: trivial path.
	if _, ok := distB[p.Start]; ok {
		res.Found = true
		res.Path = []N{p.Start}
		res.Order = []N{p.Start}
		return res, nil
	}

	meetFound := false
	var meet N

	expandF := func() {
		var next []N
		for _, node := range frontierF {
			res.Order = append(res.Order, node)
			for _, e := range p.Neighbors(node) {
				if _, ok := distF[e.To]; ok {
					continue
				}
				distF[e.To] = distF[node] + e.Weight
				parentF[e.To] = node
				if _, ok := distB[e.To]; ok {
					meet, meetFound = e.To, true
					return
				}
				next = append(next, e.To)
			}
		}
		frontierF = next
	}

	expandB := func() {
		var next []N
		for _, node := range frontierB {
			res.Order = append(res.Order, node)
			for _, e := range p.Predecessors(node) { // e.To is a predecessor of node
				if _, ok := distB[e.To]; ok {
					continue
				}
				distB[e.To] = distB[node] + e.Weight
				parentB[e.To] = node
				if _, ok := distF[e.To]; ok {
					meet, meetFound = e.To, true
					return
				}
				next = append(next, e.To)
			}
		}
		frontierB = next
	}

	for len(frontierF) > 0 && len(frontierB) > 0 && !meetFound {
		if err := p.cancelled(); err != nil {
			return res, err
		}
		if len(frontierF) <= len(frontierB) {
			expandF()
		} else {
			expandB()
		}
	}

	if !meetFound {
		return res, nil // frontiers never met: goal unreachable
	}

	// start..meet via the forward tree.
	var left []N
	for n := meet; ; {
		left = append(left, n)
		if rootF[n] {
			break
		}
		n = parentF[n]
	}
	for i, j := 0, len(left)-1; i < j; i, j = i+1, j-1 {
		left[i], left[j] = left[j], left[i]
	}

	// meet..goal via the backward tree (meet already in left).
	path := left
	if !rootB[meet] {
		for n := parentB[meet]; ; {
			path = append(path, n)
			if rootB[n] {
				break
			}
			n = parentB[n]
		}
	}

	res.Found = true
	res.Path = path
	res.Cost = distF[meet] + distB[meet]
	return res, nil
}
