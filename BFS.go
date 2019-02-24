package gotraverse

import (
	"fmt"
)
//Performs BFS on the given node until the goal node is found. Note that, it doesn't take weight and heuristic into consideration
type BFS struct {
	startNode *Node
	goalNode *Node
	capturedNode *Node
}

func (s *BFS) Search() GoalSearch {
	fmt.Println("BFS")

	queue := []*Node {s.startNode}

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
		if removedNode == s.goalNode {
			s.capturedNode = removedNode
		}
		for _, v := range removedNode.edges{
			if !v.to.isClosed {
				v.to.isClosed = true
				queue = append(queue, v.to)
			}
		}

	}
	s.capturedNode = new(Node)
	return s
}

func (s *BFS) SetIterativeNode(node *Node) GoalSearch {
	s.startNode = node
	return s
}

func (s *BFS) SetGoalNode(node *Node) GoalSearch {
	s.goalNode = node
	return s
}

func (s *BFS) GetCapturedNode() *Node {
	return s.capturedNode
}

