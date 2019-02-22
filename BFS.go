package gotraverse

import (
	"fmt"
)
//Performs BFS on the given node until the goal node is found. Note that, it doesn't take weight and heuristic into consideration
func BFS(startNode *Node, goalNode *Node) *Node {
	fmt.Println("BFS")

	queue := []*Node {startNode}

	for len(queue)>0 {
		if LOG_BFS_STEPS {
			fmt.Print("[")
			for _, n := range queue {
				fmt.Print(n.name, " ")
			}
			fmt.Println("]")
		}
		removedNode := queue[0]
		if LOG_BFS_DETAILS {
			fmt.Println("Dequeue Node:",removedNode.name)
		}
		//Remove first Node ie. Dequeue
		queue = queue[1:]
		if removedNode == goalNode {
			return goalNode
		}
		for _, v := range removedNode.edges{
			if !v.to.isClosed {
				v.to.isClosed = true
				queue = append(queue, v.to)
			}
		}

	}
	return new(Node)
}