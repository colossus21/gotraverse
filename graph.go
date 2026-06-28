package gotraverse

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Graph is an explicit, hand-built weighted directed graph with optional
// per-node heuristics. It is a convenience for the common case where the whole
// graph is known up front; call [Graph.Problem] to turn it into a [Problem] for
// the search functions. For large or infinite state spaces, build a Problem
// directly with an on-demand Neighbors function instead.
//
// Adjacency (both out- and in-edges) is kept in insertion order, so searches —
// including [Bidirectional], which relies on predecessors — are deterministic
// without requiring nodes to be ordered.
type Graph[N comparable] struct {
	has  map[N]bool
	out  map[N][]Edge[N]
	in   map[N][]Edge[N]
	heur map[N]float64
}

// New returns an empty graph.
func New[N comparable]() *Graph[N] {
	return &Graph[N]{
		has:  make(map[N]bool),
		out:  make(map[N][]Edge[N]),
		in:   make(map[N][]Edge[N]),
		heur: make(map[N]float64),
	}
}

// AddNode declares a node. Declaring the same node twice is harmless.
func (g *Graph[N]) AddNode(n N) {
	g.has[n] = true
}

// SetHeuristic declares node n and sets its heuristic value.
func (g *Graph[N]) SetHeuristic(n N, h float64) {
	g.has[n] = true
	g.heur[n] = h
}

// AddEdge adds a directed edge from -> to with the given weight. Both endpoints
// must already be declared (via [Graph.AddNode] or [Graph.SetHeuristic]).
func (g *Graph[N]) AddEdge(from, to N, weight float64) error {
	if !g.has[from] {
		return fmt.Errorf("gotraverse: edge from unknown node %v", from)
	}
	if !g.has[to] {
		return fmt.Errorf("gotraverse: edge to unknown node %v", to)
	}
	g.out[from] = append(g.out[from], Edge[N]{To: to, Weight: weight})
	g.in[to] = append(g.in[to], Edge[N]{To: from, Weight: weight})
	return nil
}

// Has reports whether a node exists in the graph.
func (g *Graph[N]) Has(n N) bool { return g.has[n] }

// Problem returns a single-goal Problem over the graph. The returned Problem
// has Neighbors, Predecessors, Heuristic and GoalNodes wired up, so every
// algorithm — including [Bidirectional] — works with it.
func (g *Graph[N]) Problem(start, goal N) Problem[N] {
	p := g.ProblemFunc(start, GoalNode(goal))
	p.GoalNodes = []N{goal}
	return p
}

// ProblemFunc returns a Problem with a custom goal predicate (e.g. to accept
// any of several goals). The result omits GoalNodes, so [Bidirectional] is not
// available; use [Graph.Problem] when you need it.
func (g *Graph[N]) ProblemFunc(start N, goal func(N) bool) Problem[N] {
	return Problem[N]{
		Start:        start,
		Goal:         goal,
		Neighbors:    func(n N) []Edge[N] { return g.out[n] },
		Predecessors: func(n N) []Edge[N] { return g.in[n] },
		Heuristic:    func(n N) float64 { return g.heur[n] },
	}
}

// Parse builds a string-keyed graph from two space-separated strings.
//
// nodes is a flat list of name/heuristic pairs, e.g. "S 8 A 8 G 0"; a heuristic
// of "inf" is parsed as +Inf. edges is a flat list of from/to/weight triples,
// e.g. "S A 3 S B 1". Every node referenced by an edge must be declared in
// nodes. Parse returns an error on malformed input or unknown endpoints.
func Parse(nodes, edges string) (*Graph[string], error) {
	g := New[string]()

	nf := strings.Fields(nodes)
	if len(nf)%2 != 0 {
		return nil, fmt.Errorf("gotraverse: nodes string has dangling token; want name/heuristic pairs")
	}
	for i := 0; i < len(nf); i += 2 {
		name := nf[i]
		if nf[i+1] == "inf" {
			g.SetHeuristic(name, math.Inf(1))
			continue
		}
		v, err := strconv.ParseFloat(nf[i+1], 64)
		if err != nil {
			return nil, fmt.Errorf("gotraverse: bad heuristic %q for node %q: %w", nf[i+1], name, err)
		}
		g.SetHeuristic(name, v)
	}

	ef := strings.Fields(edges)
	if len(ef)%3 != 0 {
		return nil, fmt.Errorf("gotraverse: edges string has dangling token; want from/to/weight triples")
	}
	for i := 0; i < len(ef); i += 3 {
		w, err := strconv.ParseFloat(ef[i+2], 64)
		if err != nil {
			return nil, fmt.Errorf("gotraverse: bad weight %q for edge %s->%s: %w", ef[i+2], ef[i], ef[i+1], err)
		}
		if err := g.AddEdge(ef[i], ef[i+1], w); err != nil {
			return nil, err
		}
	}

	return g, nil
}
