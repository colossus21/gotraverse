package gotraverse

import "fmt"

type DFS struct {
	startNode *Node
	goalNode *Node
	capturedNode *Node
}

func (s *DFS) Search() GoalSearch {
	fmt.Println("DFS:")

	stack := []*Node {s.startNode}

	for len(stack)>0 {
		if LOG_DFS_STEPS {
			fmt.Print("[")
			for _, n := range stack {
				fmt.Print(n.name, " ")
			}
			fmt.Println("]")
		}
		removedNode := stack[len(stack)-1]
		if LOG_DFS_DETAILS {
			fmt.Println("Pop Node: ", removedNode.name)
		}
		//Remove last element
		stack = stack[:len(stack)-1]
		if removedNode == s.goalNode {
			s.capturedNode = removedNode
			return s
		}
		//Get all children of all nodes and append
		for _, v := range removedNode.edges {
			stack = append(stack, v.to)
		}
	}
	s.capturedNode = new(Node)
	return s
}

func (s *DFS) SetIterativeNode(node *Node) GoalSearch {
	s.startNode = node
	return s
}

func (s *DFS) SetGoalNode(node *Node) GoalSearch {
	s.goalNode = node
	return s
}

func (s *DFS) GetCapturedNode() *Node {
	return s.capturedNode
}

