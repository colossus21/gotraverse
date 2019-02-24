package gotraverse

type Graph struct {
	strNodes string
	strEdges string
	generatedNode *Node
	nodes map[string]*Node
	goalSearch GoalSearch
}

func (g *Graph) GenerateIterativeNode(s string) {
	g.generatedNode = MakeEdges(g.strEdges, g.nodes, s)
}

func (g *Graph) GetGeneratedNode() *Node {
	return g.generatedNode
}

func (g *Graph) InitGoalSearch(algo GoalSearch, goalNode string) GoalSearch {
	algo.SetGoalNode(g.nodes[goalNode])
	algo.SetIterativeNode(g.GetGeneratedNode())
	return algo
}

func NewGraph(nodes string, edges string, startNode string) *Graph {
	g := new(Graph)
	g.strNodes = nodes
	g.strEdges = edges
	g.nodes = MakeNodes(nodes)
	g.generatedNode = MakeEdges(edges, g.nodes, startNode)
	return g
}