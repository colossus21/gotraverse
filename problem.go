// Package gotraverse implements classic graph search algorithms — BFS, DFS,
// Depth-Limited, Iterative-Deepening DFS, Bidirectional, Uniform-Cost, Greedy
// Best-First, A* and IDA* — over arbitrary graphs.
//
// Graphs are described by a [Problem], which is deliberately implicit: instead
// of materialising every node and edge, you supply a Neighbors function that
// generates successors on demand. This makes the algorithms work equally well
// on a tiny hand-built graph and on a huge or infinite state space (grids,
// puzzles, game states) without ever building the whole thing.
//
// Nodes can be any comparable type — strings, ints, or your own coordinate or
// state structs — and edge weights are float64. For the common case of an
// explicit, hand-listed graph, build one with [New] or [Parse] and call
// [Graph.Problem] to obtain a Problem.
//
//	g, _ := gotraverse.Parse(
//		"S 8 A 8 B 4 C 3 D inf E inf G 0",
//		"S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5",
//	)
//	res, _ := gotraverse.AStar(g.Problem("S", "G"))
//	fmt.Println(res.Path, res.Cost) // [S C G] 13
//
// For an implicit graph, populate a Problem directly:
//
//	p := gotraverse.Problem[Point]{
//		Start:     Point{0, 0},
//		Goal:      gotraverse.GoalNode(Point{9, 9}),
//		Neighbors: func(pt Point) []gotraverse.Edge[Point] { ... },
//		Heuristic: func(pt Point) float64 { ... },
//	}
//	res, _ := gotraverse.AStar(p)
package gotraverse

import "errors"

// Edge is a weighted connection to a node.
type Edge[N comparable] struct {
	To     N
	Weight float64
}

// Problem describes a search over a (possibly implicit) graph.
//
// Neighbors and Goal are required. Heuristic is optional and used only by the
// informed searches (Greedy, A*, IDA*); when nil it is treated as the zero
// heuristic, which degrades A* and IDA* to uniform-cost search and Greedy to an
// unguided search. Predecessors and GoalNodes are required only by
// [Bidirectional].
type Problem[N comparable] struct {
	// Start is the node the search begins from.
	Start N
	// Goal reports whether a node satisfies the goal. Use [GoalNode] for the
	// common single-target case, or any predicate to accept multiple goals.
	Goal func(N) bool
	// Neighbors returns the outgoing edges of a node, generated on demand. The
	// returned order is preserved by the algorithms, giving deterministic
	// results.
	Neighbors func(N) []Edge[N]
	// Heuristic estimates the remaining cost from a node to the goal. For A*
	// and IDA* to return optimal paths it must be admissible (never an
	// overestimate). Optional.
	Heuristic func(N) float64
	// Predecessors returns the incoming edges of a node (each Edge.To is a
	// predecessor, Edge.Weight the weight of predecessor->node). Required by
	// [Bidirectional] for its backward search.
	Predecessors func(N) []Edge[N]
	// GoalNodes lists concrete goal nodes for [Bidirectional] to seed its
	// backward search from. [Graph.Problem] sets this automatically.
	GoalNodes []N
}

// GoalNode returns a goal predicate matching a single target node.
func GoalNode[N comparable](target N) func(N) bool {
	return func(n N) bool { return n == target }
}

func (p Problem[N]) validate() error {
	if p.Neighbors == nil {
		return errors.New("gotraverse: Problem.Neighbors is nil")
	}
	if p.Goal == nil {
		return errors.New("gotraverse: Problem.Goal is nil")
	}
	return nil
}

// h returns the heuristic value for n, treating a nil Heuristic as zero.
func (p Problem[N]) h(n N) float64 {
	if p.Heuristic == nil {
		return 0
	}
	return p.Heuristic(n)
}
