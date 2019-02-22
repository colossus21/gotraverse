package gotraverse

import (
	"container/heap"
	"fmt"
)
// Performs Greedy Search on the given node until the goal node is found.
func GreedySearch(startNode *Node, goalNode *Node) *Node {
	pq := PriorityQueue{}
	heap.Init(&pq)
	heap.Push(&pq, NodeItem(startNode, startNode.heuristicValue))
	for pq.Len()>0 {
		if (LOG_GREEDY_STEPS) {
			for _, n := range pq {
				fmt.Print("[", n.value.name, " ", n.priority, "]")
			}
			fmt.Println()
		}
		removedNode := heap.Pop(&pq).(*Item)
		removedNode.value.isClosed = true
		if (LOG_GREEDY_STEPS) {
			fmt.Println(removedNode.value.name)
		}
		if (removedNode.value == goalNode) {
			return removedNode.value
		}
		for _, v := range removedNode.value.edges {
			if (LOG_GREEDY_DETAILS) {
				fmt.Println("EDGE FOUND:", v.to.name, "Heuristic:", v.to.heuristicValue, "isClosed?", v.to.isClosed)
			}
			// Conditions:
			// 1) If child is CLOSED don't add it to the queue
			// 2) If current F(N) is larger, set the F(N) and PARENT node
			// 3) For CONDITION 2, make sure the NODE WAS VISITED BEFORE, ie. Parent != Nil
			is_not_closed := !v.to.isClosed
			has_parent := v.to.parent != nil
			nodeTo := v.to
			nodeFrom := v.from

			if is_not_closed && !has_parent{
				nodeTo.parent = nodeFrom
				heap.Push(&pq, NodeItem(nodeTo, nodeTo.heuristicValue))
				if (LOG_GREEDY_DETAILS) {
					fmt.Println("EDGE ADDED:", nodeTo.name, "Heuristic:", nodeTo.heuristicValue, "Distance:", nodeTo.distanceFromStartNode, "Total:", nodeTo.GetTotalDistance(), "isClosed?", v.to.isClosed)
				}
			}
		}
	}
	panic("Node not found!!")
}

