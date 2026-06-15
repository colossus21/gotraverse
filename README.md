# GoTraverse

[![CI](https://github.com/colossus21/gotraverse/actions/workflows/ci.yml/badge.svg)](https://github.com/colossus21/gotraverse/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/colossus21/gotraverse.svg)](https://pkg.go.dev/github.com/colossus21/gotraverse)
[![Go Report Card](https://goreportcard.com/badge/github.com/colossus21/gotraverse)](https://goreportcard.com/report/github.com/colossus21/gotraverse)

A small Go library of classic graph search algorithms over a weighted directed
graph with per-node heuristics:

| Algorithm | Uses weights | Uses heuristic | Optimal |
|-----------|:------------:|:--------------:|:-------:|
| **BFS** (breadth-first)       | – | – | fewest edges |
| **DFS** (depth-first)         | – | – | no |
| **UCS** (uniform cost)        | ✓ | – | yes (min cost) |
| **Greedy** (best-first)       | – | ✓ | no |
| **A\*** | ✓ | ✓ | yes, with an admissible heuristic |

## Install

```sh
go get github.com/colossus21/gotraverse
```

```go
import "github.com/colossus21/gotraverse"
```

Requires Go 1.23+.

## Quick start

Build a graph, then run any algorithm through the `Algorithm` strategy
interface. Nodes are `name heuristic` pairs; edges are `from to weight` triples.
A heuristic of `inf` is treated as unreachable-by-heuristic.

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

	res, err := g.Search(gotraverse.AStar{}, "S", "G")
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Found) // true
	fmt.Println(res.Path)  // [S C G]
	fmt.Println(res.Cost)  // 13
	fmt.Println(res.Order) // [S B A C G]  (expansion order)
}
```

### Building a graph programmatically

```go
g := gotraverse.NewGraph()
g.AddNode("S", 8)
g.AddNode("G", 0)
g.AddEdge("S", "G", 7) // returns an error if an endpoint is undeclared
```

### Result

`Search` returns a `Result`:

```go
type Result struct {
	Algorithm string   // e.g. "A*"
	Found     bool     // was the goal reached?
	Path      []string // start..goal, inclusive (nil if not found)
	Cost      int      // total edge weight along Path
	Order     []string // nodes in the order they were expanded
}
```

### Custom algorithms

`Algorithm` is a two-method interface, so you can plug in your own strategy:

```go
type Algorithm interface {
	Name() string
	Search(g *gotraverse.Graph, start, goal string) (gotraverse.Result, error)
}
```

## Demo

The graph used above:

![graph](/img/Graph.png)

Run every algorithm against it:

```sh
go run ./example
```

```
== BFS ==
expanded : S -> A -> B -> C -> D -> E -> G
path     : S -> A -> G
cost     : 18

== DFS ==
expanded : S -> C -> G
path     : S -> C -> G
cost     : 13

== UCS ==
expanded : S -> B -> A -> D -> C -> E -> G
path     : S -> C -> G
cost     : 13

== Greedy ==
expanded : S -> C -> G
path     : S -> C -> G
cost     : 13

== A* ==
expanded : S -> B -> A -> C -> G
path     : S -> C -> G
cost     : 13
```

BFS returns `S -> A -> G` (fewest edges; `A` is declared before `C`) while the
cost-aware algorithms find the cheaper `S -> C -> G`.

## Development

```sh
go build ./...
go vet ./...
go test -race ./...
```

## License

See repository for license details.

## Acknowledgments

- Priority queue based on the [`container/heap`](https://pkg.go.dev/container/heap#example-package-PriorityQueue) example from the Go standard library.
