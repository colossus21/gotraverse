# GoTraverse

[![CI](https://github.com/colossus21/gotraverse/actions/workflows/ci.yml/badge.svg)](https://github.com/colossus21/gotraverse/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/colossus21/gotraverse.svg)](https://pkg.go.dev/github.com/colossus21/gotraverse)
[![Go Report Card](https://goreportcard.com/badge/github.com/colossus21/gotraverse)](https://goreportcard.com/report/github.com/colossus21/gotraverse)

A small, generic Go library of classic graph search algorithms. Nodes are any
`comparable` type, weights are `float64`, and graphs may be **explicit**
(hand-built) or **implicit** (neighbours generated on demand), so the same A*
runs on a five-node demo or a million-cell grid:

| Algorithm | Type | Uses weights | Uses heuristic | Optimal |
|-----------|------|:------------:|:--------------:|:-------:|
| **BFS** (breadth-first)            | uninformed | – | – | fewest edges |
| **DFS** (depth-first)              | uninformed | – | – | no |
| **Depth-Limited** (DFS to a depth) | uninformed | – | – | no |
| **IDDFS** (iterative deepening)    | uninformed | – | – | fewest edges |
| **Bidirectional** (BFS both ends)  | uninformed | – | – | fewest edges |
| **UCS** (uniform cost)             | informed   | ✓ | – | yes (min cost) |
| **Greedy** (best-first)            | informed   | – | ✓ | no |
| **A\***                            | informed   | ✓ | ✓ | yes, with an admissible heuristic |
| **IDA\*** (iterative-deepening A\*)| informed   | ✓ | ✓ | yes, with an admissible heuristic |

All algorithms share the `SearchFunc[N]` signature, so they are interchangeable;
`DepthLimited[N](limit)` is the one configured by argument, e.g.
`gotraverse.DepthLimited[string](5)`.

## Install

```sh
go get github.com/colossus21/gotraverse
```

```go
import "github.com/colossus21/gotraverse"
```

Requires Go 1.23+.

## Quick start

Nodes can be any `comparable` type and edge weights are `float64`. Searches are
plain functions with the signature `func(Problem[N]) (Result[N], error)`, so the
node type is inferred from the argument — no type parameters at the call site.

### Explicit graph

For a small, hand-listed graph, build one with `Parse` (`name heuristic` pairs
and `from to weight` triples; `inf` heuristic means unreachable-by-heuristic),
then call `Problem` to get something to search:

```go
package main

import (
	"fmt"

	"github.com/colossus21/gotraverse"
)

func main() {
	g, err := gotraverse.Parse(
		"S 8 A 8 B 4 C 3 D inf E inf G 0",
		"S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5",
	)
	if err != nil {
		panic(err)
	}

	res, _ := gotraverse.AStar(g.Problem("S", "G"))
	fmt.Println(res.Found) // true
	fmt.Println(res.Path)  // [S C G]
	fmt.Println(res.Cost)  // 13
	fmt.Println(res.Order) // [S B A C G]  (expansion order)
}
```

Build one programmatically with any node type:

```go
g := gotraverse.New[int]()
g.SetHeuristic(1, 8)
g.AddNode(2)
g.AddEdge(1, 2, 7) // returns an error if an endpoint is undeclared
res, _ := gotraverse.AStar(g.Problem(1, 2))
```

### Implicit graph (the interesting part)

You don't need to materialise the graph at all. Supply a `Neighbors` function
that generates successors on demand, and the algorithms work on grids, puzzles,
or any state space — even an unbounded one. Here's A* on a 2D grid with walls:

```go
type Point struct{ X, Y int }

p := gotraverse.Problem[Point]{
	Start: Point{0, 0},
	Goal:  gotraverse.GoalNode(Point{9, 9}),
	Neighbors: func(pt Point) []gotraverse.Edge[Point] {
		var es []gotraverse.Edge[Point]
		for _, d := range []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
			n := Point{pt.X + d.X, pt.Y + d.Y}
			if walkable(n) { // your own bounds/obstacle check
				es = append(es, gotraverse.Edge[Point]{To: n, Weight: 1})
			}
		}
		return es
	},
	Heuristic: func(pt Point) float64 { // Manhattan distance
		return math.Abs(float64(pt.X-9)) + math.Abs(float64(pt.Y-9))
	},
}

res, _ := gotraverse.AStar(p)
```

`Goal` is a predicate, so you can accept multiple goals (`func(n N) bool`).
`Heuristic` is optional; omit it and A*/IDA* degrade to uniform-cost search.
`Bidirectional` additionally needs `Predecessors` and `GoalNodes` — both are
filled in for you by `Graph.Problem`.

### Cancellation

Searching a large or unbounded space can run a long time, so every algorithm
honours a `context.Context`. Attach one with `WithContext` (like
`http.Request`); the search checks it once per node expansion and returns the
context's error if it is cancelled or times out:

```go
ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
defer cancel()

res, err := gotraverse.AStar(p.WithContext(ctx))
if errors.Is(err, context.DeadlineExceeded) {
	// gave up before finding a path
}
```

A `Problem` with no context (the default) never cancels.

### Result and strategies

```go
type Result[N comparable] struct {
	Algorithm string  // e.g. "A*"
	Found     bool    // was a goal reached?
	Path      []N     // start..goal, inclusive (nil if not found)
	Cost      float64 // total edge weight along Path
	Order     []N     // nodes in the order they were expanded
}
```

Because every algorithm shares the `SearchFunc[N]` signature, they're
interchangeable — store them in a slice and run each:

```go
for _, search := range []gotraverse.SearchFunc[string]{
	gotraverse.BFS[string], gotraverse.UCS[string], gotraverse.AStar[string],
} {
	res, _ := search(p)
	fmt.Println(res.Algorithm, res.Path)
}
```

`DepthLimited[N](limit)` returns a `SearchFunc[N]`; the rest are `SearchFunc[N]`
directly. You can plug in your own by writing a function with the same signature.

## Demo

`go run ./example` runs every algorithm on the explicit graph below, then solves
an implicit grid with A*.

![graph](/img/Graph.png)

```
# Explicit graph: S -> G
BFS            [S A G]  cost=18
DFS            [S C G]  cost=13
Depth-Limited  [S A G]  cost=18
IDDFS          [S A G]  cost=18
Bidirectional  [S A G]  cost=18
UCS            [S C G]  cost=13
Greedy         [S C G]  cost=13
A*             [S C G]  cost=13
IDA*           [S C G]  cost=13

# Implicit grid: A* from S to G
S....
*###.
***#.
.#*#.
.#**G
path length: 8 steps
```

The fewest-edge searches (BFS, IDDFS, Bidirectional) return `S -> A -> G`
(`A` is declared before `C`), while the cost-aware searches (DFS aside, which
just descends the stack) find the cheaper `S -> C -> G`. The grid path is the
optimal route A* finds around the `#` walls without the graph ever being built.

## Development

```sh
go build ./...
go vet ./...
go test -race ./...
```

## License

Released under the [MIT License](LICENSE).

## Acknowledgments

- Priority queue based on the [`container/heap`](https://pkg.go.dev/container/heap#example-package-PriorityQueue) example from the Go standard library.
