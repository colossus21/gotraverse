package gotraverse

import "fmt"
////Performs DFS on the given node until the goal node is found. Note that, it doesn't take weight and heuristic into consideration
func DFS(startNode *Node, goalNode *Node) *Node {
	fmt.Println("DFS:")

	stack := []*Node {startNode}

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
		if removedNode == goalNode {
			return removedNode
		}
		//Get all children of all nodes and append
		for _, v := range removedNode.edges {
			stack = append(stack, v.to)
		}
	}
	return new(Node)
}
