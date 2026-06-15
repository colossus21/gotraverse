// Command example runs every gotraverse search algorithm against the reference
// graph from the README and prints the expansion order, the path found and its
// cost.
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/colossus21/gotraverse"
)

func main() {
	g, err := gotraverse.Parse(
		"S 8 A 8 B 4 C 3 D inf E inf G 0",
		"S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5",
	)
	if err != nil {
		log.Fatal(err)
	}

	algos := []gotraverse.Algorithm{
		gotraverse.BFS{},
		gotraverse.DFS{},
		gotraverse.DepthLimited{Limit: 5},
		gotraverse.IterativeDeepening{},
		gotraverse.Bidirectional{},
		gotraverse.UCS{},
		gotraverse.Greedy{},
		gotraverse.AStar{},
		gotraverse.IDAStar{},
	}

	for _, algo := range algos {
		res, err := g.Search(algo, "S", "G")
		if err != nil {
			log.Fatalf("%s: %v", algo.Name(), err)
		}
		fmt.Printf("== %s ==\n", res.Algorithm)
		fmt.Printf("expanded : %s\n", strings.Join(res.Order, " -> "))
		if res.Found {
			fmt.Printf("path     : %s\n", strings.Join(res.Path, " -> "))
			fmt.Printf("cost     : %d\n\n", res.Cost)
		} else {
			fmt.Printf("path     : (goal unreachable)\n\n")
		}
	}
}
