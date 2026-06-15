package gotraverse_test

import (
	"reflect"
	"testing"

	"github.com/colossus21/gotraverse"
)

// The reference graph from the README:
//
//	S(8) ─3─ A(8) ─3─ D(inf)
//	 │ ╲      ├─7─ E(inf)
//	 1  8     └─15─ G(0)
//	 │   ╲
//	 B(4) C(3)
//	 │     │
//	 20    5
//	 ╰─→ G ←╯
const (
	refNodes = "S 8 A 8 B 4 C 3 D inf E inf G 0"
	refEdges = "S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5"
)

func refGraph(t *testing.T) *gotraverse.Graph {
	t.Helper()
	g, err := gotraverse.Parse(refNodes, refEdges)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return g
}

func TestSearchAlgorithms(t *testing.T) {
	tests := []struct {
		algo     gotraverse.Algorithm
		wantPath []string
		wantCost int
	}{
		// UCS, Greedy, A* and DFS all reach G via the cheapest S-C-G route here.
		{gotraverse.UCS{}, []string{"S", "C", "G"}, 13},
		{gotraverse.AStar{}, []string{"S", "C", "G"}, 13},
		{gotraverse.Greedy{}, []string{"S", "C", "G"}, 13},
		{gotraverse.DFS{}, []string{"S", "C", "G"}, 13},
		// BFS finds the fewest-edge path; S-A-G and S-C-G both have 2 edges,
		// and A is declared before C, so BFS discovers G via A first.
		{gotraverse.BFS{}, []string{"S", "A", "G"}, 18},
		// Recursive depth-first variants visit A's edges before C's, so they
		// reach G via A (2 edges, cost 18) — the shallowest goal.
		{gotraverse.DepthLimited{Limit: 3}, []string{"S", "A", "G"}, 18},
		{gotraverse.IterativeDeepening{}, []string{"S", "A", "G"}, 18},
		{gotraverse.Bidirectional{}, []string{"S", "A", "G"}, 18},
		// IDA* is cost-optimal like A*, so it returns the cheaper S-C-G.
		{gotraverse.IDAStar{}, []string{"S", "C", "G"}, 13},
	}

	for _, tt := range tests {
		t.Run(tt.algo.Name(), func(t *testing.T) {
			g := refGraph(t)
			res, err := g.Search(tt.algo, "S", "G")
			if err != nil {
				t.Fatalf("Search: %v", err)
			}
			if !res.Found {
				t.Fatalf("goal not found")
			}
			if !reflect.DeepEqual(res.Path, tt.wantPath) {
				t.Errorf("path = %v, want %v", res.Path, tt.wantPath)
			}
			if res.Cost != tt.wantCost {
				t.Errorf("cost = %d, want %d", res.Cost, tt.wantCost)
			}
			if res.Algorithm != tt.algo.Name() {
				t.Errorf("algorithm = %q, want %q", res.Algorithm, tt.algo.Name())
			}
			if got := res.Order[0]; got != "S" {
				t.Errorf("first expanded = %q, want start S", got)
			}
			if got := res.Order[len(res.Order)-1]; got != "G" {
				t.Errorf("last expanded = %q, want goal G", got)
			}
		})
	}
}

func TestUCSIsOptimal(t *testing.T) {
	// UCS must return the minimum-cost path even when more edges are involved.
	g := refGraph(t)
	res, err := g.Search(gotraverse.UCS{}, "S", "G")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	// S-A-G = 18, S-B-G = 21, S-C-G = 13 -> 13 is optimal.
	if res.Cost != 13 {
		t.Errorf("UCS cost = %d, want optimal 13", res.Cost)
	}
}

func TestGoalUnreachable(t *testing.T) {
	// X is declared but has no incoming edges, so nothing reaches it.
	g, err := gotraverse.Parse("S 0 A 1 X 0", "S A 1")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	for _, algo := range []gotraverse.Algorithm{
		gotraverse.BFS{}, gotraverse.DFS{}, gotraverse.UCS{}, gotraverse.Greedy{}, gotraverse.AStar{},
		gotraverse.DepthLimited{Limit: 5}, gotraverse.IterativeDeepening{},
		gotraverse.IDAStar{}, gotraverse.Bidirectional{},
	} {
		res, err := g.Search(algo, "S", "X")
		if err != nil {
			t.Fatalf("%s: %v", algo.Name(), err)
		}
		if res.Found {
			t.Errorf("%s: reported reaching unreachable goal X via %v", algo.Name(), res.Path)
		}
		if res.Path != nil {
			t.Errorf("%s: path should be nil when not found, got %v", algo.Name(), res.Path)
		}
	}
}

func TestStartEqualsGoal(t *testing.T) {
	g := refGraph(t)
	for _, algo := range []gotraverse.Algorithm{
		gotraverse.BFS{}, gotraverse.DFS{}, gotraverse.UCS{}, gotraverse.Greedy{}, gotraverse.AStar{},
		gotraverse.DepthLimited{Limit: 5}, gotraverse.IterativeDeepening{},
		gotraverse.IDAStar{}, gotraverse.Bidirectional{},
	} {
		res, err := g.Search(algo, "S", "S")
		if err != nil {
			t.Fatalf("%s: %v", algo.Name(), err)
		}
		if !res.Found || !reflect.DeepEqual(res.Path, []string{"S"}) || res.Cost != 0 {
			t.Errorf("%s: start==goal -> found=%v path=%v cost=%d, want found path=[S] cost=0",
				algo.Name(), res.Found, res.Path, res.Cost)
		}
	}
}

func TestValidationErrors(t *testing.T) {
	g := refGraph(t)
	if _, err := g.Search(gotraverse.BFS{}, "Z", "G"); err == nil {
		t.Error("expected error for unknown start node")
	}
	if _, err := g.Search(gotraverse.BFS{}, "S", "Z"); err == nil {
		t.Error("expected error for unknown goal node")
	}
}

func TestParseErrors(t *testing.T) {
	cases := map[string]struct{ nodes, edges string }{
		"dangling node token":  {"S 8 A", "S A 1"},
		"dangling edge token":  {"S 8 A 8", "S A"},
		"bad heuristic":        {"S high", ""},
		"bad weight":           {"S 8 A 8", "S A heavy"},
		"edge to unknown node": {"S 8", "S Z 1"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := gotraverse.Parse(c.nodes, c.edges); err == nil {
				t.Errorf("expected error for %s", name)
			}
		})
	}
}

func TestParseInfHeuristic(t *testing.T) {
	g, err := gotraverse.Parse("S 0 D inf", "S D 1")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if h, _ := g.Heuristic("D"); h != gotraverse.Inf {
		t.Errorf("heuristic(D) = %d, want Inf (%d)", h, gotraverse.Inf)
	}
}

func TestDepthLimitCutoff(t *testing.T) {
	g := refGraph(t)
	// G is 2 edges from S, so a limit of 1 must not reach it...
	res, err := g.Search(gotraverse.DepthLimited{Limit: 1}, "S", "G")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if res.Found {
		t.Errorf("limit 1 reached G at depth 2: %v", res.Path)
	}
	// ...but a limit of 2 must.
	res, err = g.Search(gotraverse.DepthLimited{Limit: 2}, "S", "G")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if !res.Found {
		t.Error("limit 2 failed to reach G at depth 2")
	}
}

func TestProgrammaticBuild(t *testing.T) {
	g := gotraverse.NewGraph()
	g.AddNode("S", 1)
	g.AddNode("G", 0)
	if err := g.AddEdge("S", "G", 7); err != nil {
		t.Fatalf("AddEdge: %v", err)
	}
	if err := g.AddEdge("S", "Q", 1); err == nil {
		t.Error("expected error adding edge to undeclared node")
	}
	res, err := g.Search(gotraverse.AStar{}, "S", "G")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if !res.Found || res.Cost != 7 {
		t.Errorf("found=%v cost=%d, want found cost=7", res.Found, res.Cost)
	}
}
