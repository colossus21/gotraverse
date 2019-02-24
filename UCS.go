package gotraverse

import (
	"container/heap"
	"fmt"
)
//Performs UCS on the given node until the goal node is found
type UCS struct {
	startNode *Node
	goalNode *Node
	capturedNode *Node
}

func (s *UCS) Search() GoalSearch {
	pq := PriorityQueue{}
	heap.Init(&pq)
	heap.Push(&pq, NodeItem(s.startNode, s.startNode.distanceFromStartNode))
	for pq.Len()>0 {
		if (LOG_UCS_STEPS) {
			for _, n := range pq {
				fmt.Print("[", n.value.name, " ", n.priority, "]")
			}
			fmt.Println()
		}
		removedNode := heap.Pop(&pq).(*Item)
		removedNode.value.isClosed = true
		if (LOG_UCS_STEPS) {
			fmt.Println(removedNode.value.name)
		}
		if (removedNode.value == s.goalNode) {
			s.capturedNode = removedNode.value
			return s
		}
		for _, v := range removedNode.value.edges {
			if (LOG_UCS_DETAILS) {
				fmt.Println("EDGE FOUND:", v.to.name, "isClosed?", v.to.isClosed)
			}
			// Conditions:
			// 1) If child is CLOSED don't add it to the queue
			// 2) If current F(N) is larger, set the F(N) and PARENT node
			// 3) For CONDITION 2, make sure the NODE WAS VISITED BEFORE, ie. Parent != Nil
			is_not_closed := !v.to.isClosed
			is_smallerdist := v.to.GetTotalDistance() > (v.distance + v.to.heuristicValue)
			has_parent := v.to.parent != nil
			nodeTo := v.to
			nodeFrom := v.from

			if is_not_closed {
				if !has_parent {
					nodeTo.distanceFromStartNode = v.distance + nodeFrom.distanceFromStartNode
					nodeTo.parent = nodeFrom
					heap.Push(&pq, NodeItem(nodeTo, nodeTo.distanceFromStartNode))
				} else {
					if is_smallerdist {
						nodeTo.distanceFromStartNode = v.distance + nodeFrom.distanceFromStartNode
						nodeTo.parent = nodeFrom
						heap.Push(&pq, NodeItem(nodeTo, nodeTo.distanceFromStartNode))
					}
				}
				if (LOG_UCS_DETAILS) {
					fmt.Println("EDGE ADDED:", nodeTo.name, "Distance:", nodeTo.distanceFromStartNode, "isClosed?", v.to.isClosed)
				}
			}
		}
	}
	return s
}

func (s *UCS) SetIterativeNode(node *Node) GoalSearch {
	s.startNode = node
	return s
}

func (s *UCS) SetGoalNode(node *Node) GoalSearch {
	s.goalNode = node
	return s
}

func (s *UCS) GetCapturedNode() *Node {
	return s.capturedNode
}