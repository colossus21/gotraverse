package gotraverse

import "container/heap"

// bestFirst is the shared engine for the priority-queue searches: UCS, Greedy
// and A*. They differ only in the priority assigned to a node and in whether a
// cheaper path to an already-discovered node may relax it:
//
//   - UCS:    priority = g(n),        relax = true
//   - Greedy: priority = h(n),        relax = false
//   - A*:     priority = g(n) + h(n), relax = true
//
// g(n) is the accumulated edge weight from start to n; h(n) is its heuristic.
// Ties in priority break by insertion order for stable output.
func bestFirst[N comparable](p Problem[N], name string, relax bool, prio func(n N, gCost float64) float64) (Result[N], error) {
	res := Result[N]{Algorithm: name}
	if err := p.validate(); err != nil {
		return res, err
	}

	dist := map[N]float64{p.Start: 0}
	seen := map[N]bool{p.Start: true}
	closed := map[N]bool{}
	parent := map[N]N{}

	pq := &priorityQueue[N]{}
	heap.Init(pq)
	seq := 0
	heap.Push(pq, &pqItem[N]{node: p.Start, priority: prio(p.Start, 0), seq: seq})
	seq++

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*pqItem[N])
		if closed[cur.node] {
			continue // stale entry left over from a relaxation
		}
		closed[cur.node] = true
		res.Order = append(res.Order, cur.node)

		if p.Goal(cur.node) {
			res.Found = true
			res.Path = reconstructPath(parent, p.Start, cur.node)
			res.Cost = dist[cur.node]
			return res, nil
		}

		for _, e := range p.Neighbors(cur.node) {
			if closed[e.To] {
				continue
			}
			nd := dist[cur.node] + e.Weight
			switch {
			case !seen[e.To]:
				seen[e.To] = true
				dist[e.To] = nd
				parent[e.To] = cur.node
				heap.Push(pq, &pqItem[N]{node: e.To, priority: prio(e.To, nd), seq: seq})
				seq++
			case relax && nd < dist[e.To]:
				dist[e.To] = nd
				parent[e.To] = cur.node
				heap.Push(pq, &pqItem[N]{node: e.To, priority: prio(e.To, nd), seq: seq})
				seq++
			}
		}
	}

	return res, nil
}

// UCS is uniform-cost search (Dijkstra from a single source). It expands nodes
// in order of accumulated path cost g(n) and returns a minimum-cost path,
// ignoring heuristics.
//
// Like Dijkstra's algorithm it assumes non-negative edge weights; with negative
// weights the returned cost is not guaranteed optimal.
func UCS[N comparable](p Problem[N]) (Result[N], error) {
	return bestFirst(p, "UCS", true, func(_ N, gCost float64) float64 {
		return gCost
	})
}

// Greedy is greedy best-first search. It always expands the open node with the
// smallest heuristic h(n), ignoring cost already spent. It is fast but not
// optimal: the returned path is the first one found to the goal.
func Greedy[N comparable](p Problem[N]) (Result[N], error) {
	return bestFirst(p, "Greedy", false, func(n N, _ float64) float64 {
		return p.h(n)
	})
}

// AStar is A* search. It expands nodes in order of f(n) = g(n) + h(n), the
// accumulated cost plus the heuristic estimate to the goal. With an admissible
// (non-overestimating) heuristic it returns a minimum-cost path.
//
// Like UCS it assumes non-negative edge weights.
func AStar[N comparable](p Problem[N]) (Result[N], error) {
	return bestFirst(p, "A*", true, func(n N, gCost float64) float64 {
		return gCost + p.h(n)
	})
}
