// Package gotraverse implements classic graph search algorithms — BFS, DFS,
// Uniform Cost Search, Greedy Best-First Search and A* — over a small,
// explicit, weighted directed graph with per-node heuristic values.
//
// A graph is built either programmatically with [NewGraph], [Graph.AddNode] and
// [Graph.AddEdge], or parsed from two space-separated strings with [Parse].
// Searches are run through the [Algorithm] strategy interface, e.g.
//
//	g, _ := gotraverse.Parse(
//		"S 8 A 8 B 4 C 3 D inf E inf G 0",
//		"S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5",
//	)
//	res, _ := g.Search(gotraverse.AStar{}, "S", "G")
//	fmt.Println(res.Path, res.Cost) // [S C G] 13
package gotraverse

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Inf is the heuristic value assigned to nodes declared as "inf". It is large
// enough to be treated as unreachable-by-heuristic without overflowing when
// summed with realistic edge weights.
const Inf = math.MaxInt32

// edge is a directed, weighted connection to another node.
type edge struct {
	to     string
	weight int
}

// Graph is a weighted directed graph whose nodes carry heuristic values used by
// the informed search algorithms (Greedy and A*). The zero value is not usable;
// construct one with [NewGraph] or [Parse]. A Graph is immutable once a search
// runs against it — algorithms keep their own per-search state and never mutate
// the graph, so the same Graph can be searched repeatedly and concurrently.
type Graph struct {
	heuristic map[string]int
	adj       map[string][]edge
}

// NewGraph returns an empty graph ready for [Graph.AddNode] and [Graph.AddEdge].
func NewGraph() *Graph {
	return &Graph{
		heuristic: make(map[string]int),
		adj:       make(map[string][]edge),
	}
}

// AddNode declares a node and its heuristic value, overwriting any previous
// declaration of the same name.
func (g *Graph) AddNode(name string, heuristic int) {
	g.heuristic[name] = heuristic
	if _, ok := g.adj[name]; !ok {
		g.adj[name] = nil
	}
}

// AddEdge adds a directed edge from -> to with the given weight. Both endpoints
// must already be declared with [Graph.AddNode]; otherwise an error is returned.
func (g *Graph) AddEdge(from, to string, weight int) error {
	if _, ok := g.heuristic[from]; !ok {
		return fmt.Errorf("gotraverse: edge from unknown node %q", from)
	}
	if _, ok := g.heuristic[to]; !ok {
		return fmt.Errorf("gotraverse: edge to unknown node %q", to)
	}
	g.adj[from] = append(g.adj[from], edge{to: to, weight: weight})
	return nil
}

// Has reports whether a node with the given name exists in the graph.
func (g *Graph) Has(name string) bool {
	_, ok := g.heuristic[name]
	return ok
}

// Heuristic returns the heuristic value declared for name, and whether it exists.
func (g *Graph) Heuristic(name string) (int, bool) {
	h, ok := g.heuristic[name]
	return h, ok
}

// Parse builds a graph from two space-separated strings.
//
// nodes is a flat list of name/heuristic pairs, e.g. "S 8 A 8 G 0". A heuristic
// of "inf" is parsed as [Inf]. edges is a flat list of from/to/weight triples,
// e.g. "S A 3 S B 1". Every node referenced by an edge must be declared in
// nodes. Parse returns an error on malformed input or unknown edge endpoints.
func Parse(nodes, edges string) (*Graph, error) {
	g := NewGraph()

	nf := strings.Fields(nodes)
	if len(nf)%2 != 0 {
		return nil, fmt.Errorf("gotraverse: nodes string has dangling token; want name/heuristic pairs")
	}
	for i := 0; i < len(nf); i += 2 {
		name := nf[i]
		var h int
		if nf[i+1] == "inf" {
			h = Inf
		} else {
			v, err := strconv.Atoi(nf[i+1])
			if err != nil {
				return nil, fmt.Errorf("gotraverse: bad heuristic %q for node %q: %w", nf[i+1], name, err)
			}
			h = v
		}
		g.AddNode(name, h)
	}

	ef := strings.Fields(edges)
	if len(ef)%3 != 0 {
		return nil, fmt.Errorf("gotraverse: edges string has dangling token; want from/to/weight triples")
	}
	for i := 0; i < len(ef); i += 3 {
		w, err := strconv.Atoi(ef[i+2])
		if err != nil {
			return nil, fmt.Errorf("gotraverse: bad weight %q for edge %s->%s: %w", ef[i+2], ef[i], ef[i+1], err)
		}
		if err := g.AddEdge(ef[i], ef[i+1], w); err != nil {
			return nil, err
		}
	}

	return g, nil
}
