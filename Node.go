package gotraverse

import "fmt"
//Node
type Node struct {
	parent *Node
	edges []*Edge
	name string
	heuristicValue int
	distanceFromStartNode int
	isClosed bool
}
// get total distance: distance from start node + heuristic value
func (n *Node) GetTotalDistance() int{
	// g(n) + h(n)
	return n.distanceFromStartNode + n.heuristicValue
}
// traverse a node and its children
func (n *Node) Traverse(msg... string) {
	newNode := new(Node)
	newNode.name = n.name
	newNode.parent = n.parent
	newNode.heuristicValue = n.heuristicValue
	newNode.distanceFromStartNode = n.distanceFromStartNode

	g := "G:"
	h := "H:"
	f := "F:"
	node := "Node:"

	if (msg != nil && len(msg)>3) {
		node = msg[0]
		g = msg[1]
		h = msg[2]
		f = msg[3]
	}

	for newNode != nil {
		defer fmt.Println(node, newNode.name, g, newNode.distanceFromStartNode, h, newNode.heuristicValue, f, newNode.heuristicValue)
		newNode = newNode.parent
	}
}
// Edge
type Edge struct {
	from *Node
	to *Node
	distance int
}