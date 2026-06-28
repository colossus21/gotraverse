package gotraverse

// SearchFunc is any search algorithm: it takes a [Problem] and returns a
// [Result]. The package's algorithms ([BFS], [DFS], [UCS], [Greedy], [AStar],
// [IDAStar], [IDDFS], [Bidirectional]) are SearchFuncs, and [DepthLimited]
// returns one. Because they share this signature you can store them in a slice
// and run each against the same problem:
//
//	for _, search := range []gotraverse.SearchFunc[string]{
//		gotraverse.BFS[string], gotraverse.AStar[string],
//	} {
//		res, _ := search(p)
//		fmt.Println(res.Algorithm, res.Path)
//	}
type SearchFunc[N comparable] func(Problem[N]) (Result[N], error)
