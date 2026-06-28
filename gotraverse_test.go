package gotraverse_test

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/colossus21/gotraverse"
)

const (
	refNodes = "S 8 A 8 B 4 C 3 D inf E inf G 0"
	refEdges = "S A 3 S B 1 S C 8 A D 3 A E 7 A G 15 B G 20 C G 5"
)

func refProblem(t *testing.T) gotraverse.Problem[string] {
	t.Helper()
	g, err := gotraverse.Parse(refNodes, refEdges)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return g.Problem("S", "G")
}

func TestSearchAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		search   gotraverse.SearchFunc[string]
		wantPath []string
		wantCost float64
	}{
		// Fewest-edge searches reach G via A (S-A-G and S-C-G both have 2 edges,
		// and A is declared before C).
		{"BFS", gotraverse.BFS[string], []string{"S", "A", "G"}, 18},
		{"IDDFS", gotraverse.IDDFS[string], []string{"S", "A", "G"}, 18},
		{"Bidirectional", gotraverse.Bidirectional[string], []string{"S", "A", "G"}, 18},
		{"Depth-Limited", gotraverse.DepthLimited[string](3), []string{"S", "A", "G"}, 18},
		// DFS's stack pops C before A, so it descends S-C-G.
		{"DFS", gotraverse.DFS[string], []string{"S", "C", "G"}, 13},
		// Cost-aware searches find the cheaper S-C-G.
		{"UCS", gotraverse.UCS[string], []string{"S", "C", "G"}, 13},
		{"Greedy", gotraverse.Greedy[string], []string{"S", "C", "G"}, 13},
		{"A*", gotraverse.AStar[string], []string{"S", "C", "G"}, 13},
		{"IDA*", gotraverse.IDAStar[string], []string{"S", "C", "G"}, 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.search(refProblem(t))
			if err != nil {
				t.Fatalf("search: %v", err)
			}
			if !res.Found {
				t.Fatal("goal not found")
			}
			if !reflect.DeepEqual(res.Path, tt.wantPath) {
				t.Errorf("path = %v, want %v", res.Path, tt.wantPath)
			}
			if res.Cost != tt.wantCost {
				t.Errorf("cost = %v, want %v", res.Cost, tt.wantCost)
			}
			if res.Algorithm != tt.name {
				t.Errorf("algorithm = %q, want %q", res.Algorithm, tt.name)
			}
			if len(res.Order) == 0 || res.Order[0] != "S" {
				t.Errorf("first expanded = %v, want start S", res.Order)
			}
		})
	}
}

func TestUCSIsOptimal(t *testing.T) {
	// S-A-G = 18, S-B-G = 21, S-C-G = 13 -> 13 is optimal.
	res, err := gotraverse.UCS(refProblem(t))
	if err != nil {
		t.Fatalf("UCS: %v", err)
	}
	if res.Cost != 13 {
		t.Errorf("UCS cost = %v, want optimal 13", res.Cost)
	}
}

func TestGoalUnreachable(t *testing.T) {
	// X is declared but has no incoming edges, so nothing reaches it.
	g, err := gotraverse.Parse("S 0 A 1 X 0", "S A 1")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	p := g.Problem("S", "X")
	for name, search := range allAlgorithms() {
		res, err := search(p)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if res.Found {
			t.Errorf("%s: reported reaching unreachable goal X via %v", name, res.Path)
		}
		if res.Path != nil {
			t.Errorf("%s: path should be nil when not found, got %v", name, res.Path)
		}
	}
}

func TestStartEqualsGoal(t *testing.T) {
	g, err := gotraverse.Parse(refNodes, refEdges)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	p := g.Problem("S", "S")
	for name, search := range allAlgorithms() {
		res, err := search(p)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !res.Found || !reflect.DeepEqual(res.Path, []string{"S"}) || res.Cost != 0 {
			t.Errorf("%s: start==goal -> found=%v path=%v cost=%v, want found [S] 0",
				name, res.Found, res.Path, res.Cost)
		}
	}
}

func TestDepthLimitCutoff(t *testing.T) {
	p := refProblem(t)
	if res, _ := gotraverse.DepthLimited[string](1)(p); res.Found {
		t.Errorf("limit 1 reached G at depth 2: %v", res.Path)
	}
	if res, _ := gotraverse.DepthLimited[string](2)(p); !res.Found {
		t.Error("limit 2 failed to reach G at depth 2")
	}
}

func TestImplicitGraph(t *testing.T) {
	// An implicit, never-materialised graph: from n you may step +1 or +3, each
	// costing 1. No Graph is built — neighbours are generated on demand.
	const goal = 5
	p := gotraverse.Problem[int]{
		Start: 0,
		Goal:  gotraverse.GoalNode(goal),
		Neighbors: func(n int) []gotraverse.Edge[int] {
			if n >= goal {
				return nil
			}
			return []gotraverse.Edge[int]{{To: n + 1, Weight: 1}, {To: n + 3, Weight: 1}}
		},
		Heuristic: func(n int) float64 {
			if d := goal - n; d > 0 {
				return float64(d) / 3 // admissible: at most +3 per unit cost
			}
			return 0
		},
	}

	res, err := gotraverse.AStar(p)
	if err != nil {
		t.Fatalf("AStar: %v", err)
	}
	if !res.Found {
		t.Fatal("did not reach goal")
	}
	if res.Cost != 3 { // 5 = 3+1+1 or 1+1+3 -> 3 edges minimum
		t.Errorf("cost = %v, want 3", res.Cost)
	}
	if res.Path[0] != 0 || res.Path[len(res.Path)-1] != goal {
		t.Errorf("path = %v, want 0..%d", res.Path, goal)
	}
}

func TestProblemValidation(t *testing.T) {
	// nil Neighbors
	if _, err := gotraverse.BFS(gotraverse.Problem[string]{Goal: gotraverse.GoalNode("x")}); err == nil {
		t.Error("expected error for nil Neighbors")
	}
	// nil Goal
	if _, err := gotraverse.BFS(gotraverse.Problem[string]{
		Neighbors: func(string) []gotraverse.Edge[string] { return nil },
	}); err == nil {
		t.Error("expected error for nil Goal")
	}
}

func TestBidirectionalRequirements(t *testing.T) {
	// ProblemFunc omits GoalNodes, so Bidirectional must refuse it.
	g, err := gotraverse.Parse(refNodes, refEdges)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	p := g.ProblemFunc("S", gotraverse.GoalNode("G"))
	if _, err := gotraverse.Bidirectional(p); err == nil {
		t.Error("expected error: Bidirectional needs GoalNodes")
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
	if h := g.Problem("S", "D").Heuristic("D"); !math.IsInf(h, 1) {
		t.Errorf("heuristic(D) = %v, want +Inf", h)
	}
}

func TestProgrammaticBuild(t *testing.T) {
	g := gotraverse.New[int]()
	g.AddNode(1)
	g.AddNode(2)
	if err := g.AddEdge(1, 2, 7); err != nil {
		t.Fatalf("AddEdge: %v", err)
	}
	if err := g.AddEdge(1, 99, 1); err == nil {
		t.Error("expected error adding edge to undeclared node")
	}
	res, err := gotraverse.AStar(g.Problem(1, 2))
	if err != nil {
		t.Fatalf("AStar: %v", err)
	}
	if !res.Found || res.Cost != 7 {
		t.Errorf("found=%v cost=%v, want found 7", res.Found, res.Cost)
	}
}

func TestContextCancelled(t *testing.T) {
	g, err := gotraverse.Parse(refNodes, refEdges)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already done before any search runs

	for name, search := range allAlgorithms() {
		p := g.Problem("S", "G").WithContext(ctx)
		res, err := search(p)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("%s: err = %v, want context.Canceled", name, err)
		}
		if res.Found {
			t.Errorf("%s: found goal despite cancellation", name)
		}
	}
}

func TestContextTimeoutStopsInfiniteSearch(t *testing.T) {
	// An infinite graph (n -> n+1 forever) with a goal that never matches:
	// only the context can stop these searches. If cancellation were broken
	// the test would hang and fail via the test timeout.
	mkProblem := func(ctx context.Context) gotraverse.Problem[int] {
		return gotraverse.Problem[int]{
			Start: 0,
			Goal:  func(int) bool { return false },
			Neighbors: func(n int) []gotraverse.Edge[int] {
				return []gotraverse.Edge[int]{{To: n + 1, Weight: 1}}
			},
			Context: ctx,
		}
	}

	cases := map[string]gotraverse.SearchFunc[int]{
		"BFS":           gotraverse.BFS[int],
		"DFS":           gotraverse.DFS[int],
		"UCS":           gotraverse.UCS[int],
		"Greedy":        gotraverse.Greedy[int],
		"A*":            gotraverse.AStar[int],
		"IDDFS":         gotraverse.IDDFS[int],
		"IDA*":          gotraverse.IDAStar[int],
		"Depth-Limited": gotraverse.DepthLimited[int](1 << 30),
	}
	for name, search := range cases {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		res, err := search(mkProblem(ctx))
		cancel()
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("%s: err = %v, want context.DeadlineExceeded", name, err)
		}
		if res.Found {
			t.Errorf("%s: found a goal on an infinite graph", name)
		}
	}
}

// allAlgorithms returns every search keyed by name, for shared edge-case tests.
func allAlgorithms() map[string]gotraverse.SearchFunc[string] {
	return map[string]gotraverse.SearchFunc[string]{
		"BFS":           gotraverse.BFS[string],
		"DFS":           gotraverse.DFS[string],
		"Depth-Limited": gotraverse.DepthLimited[string](10),
		"IDDFS":         gotraverse.IDDFS[string],
		"Bidirectional": gotraverse.Bidirectional[string],
		"UCS":           gotraverse.UCS[string],
		"Greedy":        gotraverse.Greedy[string],
		"A*":            gotraverse.AStar[string],
		"IDA*":          gotraverse.IDAStar[string],
	}
}
