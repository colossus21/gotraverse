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
// g(n) is the accumulated edge weight from start to n; h(n) is n's heuristic.
// prio receives both the node name and its current g so each strategy can pick
// what it needs. Ties in priority break by insertion order for stable output.
func bestFirst(g *Graph, start, goal, name string, relax bool, prio func(node string, gCost int) int) (Result, error) {
	res := Result{Algorithm: name}
	if err := g.validate(start, goal); err != nil {
		return res, err
	}

	dist := map[string]int{start: 0}
	seen := map[string]bool{start: true}
	closed := map[string]bool{}
	parent := map[string]string{}

	pq := &priorityQueue{}
	heap.Init(pq)
	seq := 0

	heap.Push(pq, &pqItem{node: start, priority: prio(start, 0), seq: seq})
	seq++

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*pqItem)
		if closed[cur.node] {
			continue // stale entry left over from a relaxation
		}
		closed[cur.node] = true
		res.Order = append(res.Order, cur.node)

		if cur.node == goal {
			res.Found = true
			res.Path = reconstructPath(parent, start, goal)
			res.Cost = dist[goal]
			return res, nil
		}

		for _, e := range g.adj[cur.node] {
			if closed[e.to] {
				continue
			}
			nd := dist[cur.node] + e.weight
			switch {
			case !seen[e.to]:
				seen[e.to] = true
				dist[e.to] = nd
				parent[e.to] = cur.node
				heap.Push(pq, &pqItem{node: e.to, priority: prio(e.to, nd), seq: seq})
				seq++
			case relax && nd < dist[e.to]:
				dist[e.to] = nd
				parent[e.to] = cur.node
				heap.Push(pq, &pqItem{node: e.to, priority: prio(e.to, nd), seq: seq})
				seq++
			}
		}
	}

	return res, nil
}
