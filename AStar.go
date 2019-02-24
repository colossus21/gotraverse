package gotraverse

import (
	"container/heap"
	"fmt"
)
//Performs A-Star on the given node until the goal node is found
type AStar struct {
	startNode *Node
	goalNode *Node
	capturedNode *Node
}

func (s *AStar) Search() GoalSearch {
	pq := PriorityQueue{}
	heap.Init(&pq)
	heap.Push(&pq, NodeItem(s.startNode, s.startNode.heuristicValue))
	for pq.Len()>0 {
		if (LOG_A_STAR_STEPS) {
			for _, n := range pq {
				fmt.Print("[", n.value.name, " ", n.priority, "]")
			}
			fmt.Println()
		}
		removedNode := heap.Pop(&pq).(*Item)
		removedNode.value.isClosed = true
		if (LOG_A_STAR_STEPS) {
			fmt.Println(removedNode.value.name)
		}
		if (removedNode.value == s.goalNode) {
			s.capturedNode = removedNode.value
			return s
		}
		for _, v := range removedNode.value.edges {
			if (LOG_A_STAR_DETAILS) {
				fmt.Println("EDGE FOUND:", v.to.name, "Heuristic:", v.to.heuristicValue, "isClosed?", v.to.isClosed)
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
					heap.Push(&pq, NodeItem(nodeTo, nodeTo.distanceFromStartNode + nodeTo.heuristicValue))
				} else {
					if is_smallerdist {
						nodeTo.distanceFromStartNode = v.distance + nodeFrom.distanceFromStartNode
						nodeTo.parent = nodeFrom
						heap.Push(&pq, NodeItem(nodeTo, nodeTo.distanceFromStartNode + nodeTo.heuristicValue))
					}
				}
				if (LOG_A_STAR_DETAILS) {
					fmt.Println("EDGE ADDED:", nodeTo.name, "Heuristic:", nodeTo.heuristicValue, "Distance:", nodeTo.distanceFromStartNode, "Total:", nodeTo.GetTotalDistance(), "isClosed?", v.to.isClosed)
				}
			}
		}
	}
	return s
}

func (s *AStar) SetIterativeNode(node *Node) GoalSearch {
	s.startNode = node
	return s
}

func (s *AStar) SetGoalNode(node *Node) GoalSearch {
	s.goalNode = node
	return s
}

func (s *AStar) GetCapturedNode() *Node {
	return s.capturedNode
}